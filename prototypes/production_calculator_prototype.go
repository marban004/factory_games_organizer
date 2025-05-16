package prototypes

import (
	"container/list"
	"fmt"

	"gorm.io/gorm"
)

type BestrecipeResult struct {
	ID             uint
	RecipeName     string
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
	RecipeName                 string
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

// If using default recipes pass empty string array
func Calculate(desiredResourceName string, desiredRate float32, names []string, db *gorm.DB) {
	ProductionTreeNodes := make(map[int]*ProductionTreeNode)
	ExcessResources := list.New()
	findAndComputeBestrecipeForResource(desiredResourceName, desiredRate, names, ProductionTreeNodes, ExcessResources, db)
	for _, Node := range ProductionTreeNodes {
		fmt.Println("Making recipe:", Node.RecipeName, Node.MachineName, Node.MachineNumber)
		fmt.Println("Node number:", Node.NodeId, "Source nodes: ")
		if Node.SourceNodes.Len() != 0 {
			for e := Node.SourceNodes.Front(); e != nil; e = e.Next() {
				fmt.Println("\t", e.Value)
			}
		}
		fmt.Println("------------------------------------------")
	}
	for e := ExcessResources.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func findBestrecipe(desiredResourceName string, names []string, db *gorm.DB, bestrecipe *BestrecipeResult) {
	var query string = `SELECT rcp.id, rcp.name AS recipe_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate 
							FROM recipes rcp
							JOIN recipes_outputs ro ON rcp.id = ro.recipes_id
							JOIN resources r ON ro.resources_id = r.id
							JOIN machines_recipes mr ON rcp.id = mr.recipes_id
							JOIN machines m ON mr.machines_id = m.id
							WHERE r.name = '` + desiredResourceName + `'`
	if len(names) == 0 {
		query += ` AND rcp.default_choice = TRUE `
	} else {
		query += " AND (rcp.default_choice = TRUE OR rcp.name IN ("
		for i, name := range names {
			if i != 0 {
				query += ","
			}
			query += "'" + name + "'"
		}
		query += ")) "
	}
	query += `ORDER BY rcp.default_choice, rate DESC
				LIMIT 1;`
	db.Raw(query).Scan(bestrecipe)
	//fmt.Println("selected recipe", bestrecipe.RecipeName)
}

func getRequiredResources(id uint, requiredResources *list.List, db *gorm.DB) {
	var resource RequiredResourceResult
	var query string = `SELECT r.id, r.name AS resource_name, ri.amount AS amount_required, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ri.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipes rcp
			JOIN recipes_inputs ri ON rcp.id = ri.recipes_id
			JOIN resources r ON ri.resources_id = r.id
			JOIN machines_recipes mr ON rcp.id = mr.recipes_id
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
	var query string = `SELECT r.id, r.name AS resource_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipes rcp
			JOIN recipes_outputs ro ON rcp.id = ro.recipes_id
			JOIN resources r ON ro.resources_id = r.id
			JOIN machines_recipes mr ON rcp.id = mr.recipes_id
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

func findAndComputeBestrecipeForResource(desiredResourceName string, desiredRate float32, names []string, ProductionTreeNodes map[int]*ProductionTreeNode, ExcessResources *list.List, db *gorm.DB) int {
	//Finding best recipe
	//fmt.Println(desiredResourceName)
	var bestrecipe BestrecipeResult
	var machinesRequired float32
	var NewNode ProductionTreeNode = ProductionTreeNode{RequiredResourcesPerSecond: make(map[string]float32), ProducedResourcesPerSecond: make(map[string]float32)}
	findBestrecipe(desiredResourceName, names, db, &bestrecipe)
	machinesRequired = desiredRate / bestrecipe.Rate
	//Found best recipe for resource
	//Finding required resources for best recipe
	requiredResources := list.New()
	producedResources := list.New()
	getRequiredResources(bestrecipe.ID, requiredResources, db)
	getProducedResources(bestrecipe.ID, producedResources, db)
	//fmt.Println(ExcessResources.Front())
	for e := requiredResources.Front(); e != nil; e = e.Next() {
		//fmt.Println(requiredResources.Len(), ExcessResources.Front())
		inserted := false
		//fmt.Println(e.Value.(RequiredResourceResult).ResourceName)
		for ei := ExcessResources.Front(); ei != nil; ei = ei.Next() {
			if ei.Value.(ResourceSource).ExcessResourceName == e.Value.(RequiredResourceResult).ResourceName {
				if ei.Value.(ResourceSource).ExcessProducedResourcePerSecond < e.Value.(RequiredResourceResult).Rate*machinesRequired {
					NewNode.RequiredResourcesPerSecond[e.Value.(RequiredResourceResult).ResourceName] = (e.Value.(RequiredResourceResult).Rate * machinesRequired) - ei.Value.(ResourceSource).ExcessProducedResourcePerSecond
				} else if ei.Value.(ResourceSource).ExcessProducedResourcePerSecond >= e.Value.(RequiredResourceResult).Rate*machinesRequired {
					newEiElement := ResourceSource{NodeId: ei.Value.(ResourceSource).NodeId, ExcessResourceName: ei.Value.(ResourceSource).ExcessResourceName, ExcessProducedResourcePerSecond: ei.Value.(ResourceSource).ExcessProducedResourcePerSecond - e.Value.(RequiredResourceResult).Rate*machinesRequired}
					if newEiElement.ExcessProducedResourcePerSecond > 0 {
						ExcessResources.InsertBefore(newEiElement, ei)
					}
				}
				ExcessResources.Remove(ei)
				inserted = true
				NewNode.SourceNodes.PushBack(ei.Value.(ResourceSource).NodeId)
			}
		}

		if !inserted {
			NewNode.RequiredResourcesPerSecond[e.Value.(RequiredResourceResult).ResourceName] = e.Value.(RequiredResourceResult).Rate * machinesRequired
		}
	}
	for e := producedResources.Front(); e != nil; e = e.Next() {
		//fmt.Println("produced resources list:", producedResources, e.Value.(ProducedResourceResult).ResourceName)
		NewNode.ProducedResourcesPerSecond[e.Value.(ProducedResourceResult).ResourceName] = e.Value.(ProducedResourceResult).Rate * machinesRequired
		if e.Value.(ProducedResourceResult).ResourceName != desiredResourceName {
			var excessResource ResourceSource
			excessResource.NodeId = len(ProductionTreeNodes)
			excessResource.ExcessResourceName = e.Value.(ProducedResourceResult).ResourceName
			excessResource.ExcessProducedResourcePerSecond = e.Value.(ProducedResourceResult).Rate * machinesRequired
			ExcessResources.PushBack(excessResource)
		}
	}
	NewNode.MachineName = bestrecipe.MachineName
	NewNode.MachineNumber = machinesRequired
	NewNode.RecipeName = bestrecipe.RecipeName
	NewNode.NodeId = len(ProductionTreeNodes)
	ProductionTreeNodes[len(ProductionTreeNodes)] = &NewNode
	//fmt.Println(NewNode.recipeName)
	for resourceName, requiredAmount := range NewNode.RequiredResourcesPerSecond {
		if requiredAmount > 0 {
			sourceNode := findAndComputeBestrecipeForResource(resourceName, requiredAmount, names, ProductionTreeNodes, ExcessResources, db)
			NewNode.SourceNodes.PushBack(sourceNode)
		}
	}
	return NewNode.NodeId
}
