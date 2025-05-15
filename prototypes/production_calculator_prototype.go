package prototypes

import (
	"container/list"
	"fmt"

	"gorm.io/gorm"
)

type BestRecipieResult struct {
	ID             uint
	RecipieName    string
	AmountProduced uint
	ProductionTime uint
	MachineName    string
	MachineSpeed   uint
	Rate           float32
}

type RequiredResourceResult struct {
	ID             uint
	ResourceName   string
	AmountRequired float32
	ProductionTime float32
	MachineName    string
	MachineSpeed   uint
	Rate           float32
}

type ProducedResourceResult struct {
	ID             uint
	ResourceName   string
	AmountProduced float32
	ProductionTime float32
	MachineName    string
	MachineSpeed   uint
	Rate           float32
}

type ProductionTreeNode struct {
	NodeId                     int
	RecipieName                string
	MachineName                string
	MachineNumber              float32
	RequiredResourcesPerSecond map[string]float32
	ProducedResourcesPerSecond map[string]float32
	SourceNodes                list.List
}

type ResourceSource struct {
	NodeId                          int
	ExcessResourceName              string
	ExcessProducedResourcePerSecond float32
}

// If using default recipies pass empty string array
func Calculate(resourceName string, desiredRate float32, names []string, db *gorm.DB) {
	ProductionTreeNodes := make(map[int]*ProductionTreeNode)
	ExcessResources := list.New()
	findAndComputeBestRecipieForResource(resourceName, desiredRate, names, ProductionTreeNodes, ExcessResources, db)
	for _, Node := range ProductionTreeNodes {
		fmt.Println("Making recipie:", Node.RecipieName, Node.MachineName, Node.MachineNumber)
		fmt.Println("Node number:", Node.NodeId, "Source nodes: ")
		if Node.SourceNodes.Len() != 0 {
			for e := Node.SourceNodes.Front(); e != nil; e = e.Next() {
				fmt.Println("\t", e.Value)
			}
		}
		fmt.Println("------------------------------------------")
	}
}

func findBestRecipie(resourceName string, names []string, db *gorm.DB, bestRecipie *BestRecipieResult) {
	var query string = `SELECT rcp.id, rcp.name AS recipie_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipies rcp
							JOIN recipies_outputs ro ON rcp.id = ro.recipies_id
							JOIN resources r ON ro.resources_id = r.id
							JOIN machines_recipies mr ON rcp.id = mr.recipies_id
							JOIN machines m ON mr.machines_id = m.id
							WHERE r.name = '` + resourceName + `' AND `
	if len(names) == 0 {
		query += "rcp.default_choice = TRUE "
	} else {
		query += "rcp.name IN ("
		for i, name := range names {
			if i != 0 {
				query += ","
			}
			query += "'" + name + "'"
		}
		query += ") "
	}
	query += `ORDER BY rate
				LIMIT 1;`
	db.Raw(query).Scan(bestRecipie)
	//fmt.Println(bestRecipie.ID)
}

func getRequiredResources(id uint, requiredResources *list.List, db *gorm.DB) {
	var resource RequiredResourceResult
	var query string = `SELECT r.id, r.name AS resource_name, ri.amount AS amount_required, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ri.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipies rcp
			JOIN recipies_inputs ri ON rcp.id = ri.recipies_id
			JOIN resources r ON ri.resources_id = r.id
			JOIN machines_recipies mr ON rcp.id = mr.recipies_id
			JOIN machines m ON mr.machines_id = m.id
			WHERE rcp.id = ?;`
	rows, err := db.Raw(query, id).Rows()
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		db.ScanRows(rows, &resource)
		requiredResources.PushBack(resource)
	}
}

func getProducedResources(id uint, producedResources *list.List, db *gorm.DB) {
	var resource ProducedResourceResult
	var query string = `SELECT r.id, r.name AS resource_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipies rcp
			JOIN recipies_outputs ro ON rcp.id = ro.recipies_id
			JOIN resources r ON ro.resources_id = r.id
			JOIN machines_recipies mr ON rcp.id = mr.recipies_id
			JOIN machines m ON mr.machines_id = m.id
			WHERE rcp.id = ?;`
	rows, err := db.Raw(query, id).Rows()
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		db.ScanRows(rows, &resource)
		producedResources.PushBack(resource)
	}
}

func findAndComputeBestRecipieForResource(resourceName string, desiredRate float32, names []string, ProductionTreeNodes map[int]*ProductionTreeNode, ExcessResources *list.List, db *gorm.DB) int {
	//Finding best recipie
	var bestRecipie BestRecipieResult
	var machinesRequired float32
	var NewNode ProductionTreeNode = ProductionTreeNode{RequiredResourcesPerSecond: make(map[string]float32), ProducedResourcesPerSecond: make(map[string]float32)}
	findBestRecipie(resourceName, names, db, &bestRecipie)
	machinesRequired = desiredRate / bestRecipie.Rate
	//fmt.Println(machinesRequired)
	//Found best recipie for resource
	//Finding required resources for best recipie
	requiredResources := list.New()
	producedResources := list.New()
	getRequiredResources(bestRecipie.ID, requiredResources, db)
	getProducedResources(bestRecipie.ID, producedResources, db)
	for e := requiredResources.Front(); e != nil; e = e.Next() {
		inserted := false
		//fmt.Println(e.Value.(RequiredResourceResult).ResourceName, e.Value.(RequiredResourceResult).Rate*machinesRequired)
		for ei := ExcessResources.Front(); ei != nil; ei = ei.Next() {
			if ei.Value.(ResourceSource).ExcessResourceName == e.Value.(RequiredResourceResult).ResourceName {
				if ei.Value.(ResourceSource).ExcessProducedResourcePerSecond < e.Value.(RequiredResourceResult).Rate*machinesRequired {
					NewNode.RequiredResourcesPerSecond[e.Value.(RequiredResourceResult).ResourceName] = (e.Value.(RequiredResourceResult).Rate * machinesRequired) - ei.Value.(ResourceSource).ExcessProducedResourcePerSecond
					ExcessResources.Remove(ei)
					inserted = true
				} else if ei.Value.(ResourceSource).ExcessProducedResourcePerSecond >= e.Value.(RequiredResourceResult).Rate*machinesRequired {
					newEiElement := ResourceSource{NodeId: ei.Value.(ResourceSource).NodeId, ExcessResourceName: ei.Value.(ResourceSource).ExcessResourceName, ExcessProducedResourcePerSecond: ei.Value.(ResourceSource).ExcessProducedResourcePerSecond - e.Value.(RequiredResourceResult).Rate*machinesRequired}
					if newEiElement.ExcessProducedResourcePerSecond > 0 {
						ExcessResources.InsertBefore(newEiElement, ei)
					}
					ExcessResources.Remove(ei)
					inserted = true
				}
			}
		}
		if !inserted {
			NewNode.RequiredResourcesPerSecond[e.Value.(RequiredResourceResult).ResourceName] = e.Value.(RequiredResourceResult).Rate * machinesRequired
		}
	}
	for e := producedResources.Front(); e != nil; e = e.Next() {
		//fmt.Println(e.Value.(ProducedResourceResult).ResourceName, e.Value.(ProducedResourceResult).AmountProduced*machinesRequired)
		NewNode.ProducedResourcesPerSecond[e.Value.(ProducedResourceResult).ResourceName] = e.Value.(ProducedResourceResult).Rate * machinesRequired
		if e.Value.(ProducedResourceResult).ResourceName != resourceName {
			var excessResource ResourceSource
			excessResource.NodeId = len(ProductionTreeNodes)
			excessResource.ExcessResourceName = e.Value.(ProducedResourceResult).ResourceName
			excessResource.ExcessProducedResourcePerSecond = e.Value.(ProducedResourceResult).Rate * machinesRequired
			ExcessResources.PushBack(excessResource)
		}
	}
	NewNode.MachineName = bestRecipie.MachineName
	NewNode.MachineNumber = machinesRequired
	NewNode.RecipieName = bestRecipie.RecipieName
	NewNode.NodeId = len(ProductionTreeNodes)
	ProductionTreeNodes[len(ProductionTreeNodes)] = &NewNode
	//fmt.Println(NewNode.RecipieName)
	for resourceName, requiredAmount := range NewNode.RequiredResourcesPerSecond {
		sourceNode := findAndComputeBestRecipieForResource(resourceName, requiredAmount, names, ProductionTreeNodes, ExcessResources, db)
		NewNode.SourceNodes.PushBack(sourceNode)
	}
	return NewNode.NodeId
}
