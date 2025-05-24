package microservicelogiccalculator

import (
	"container/list"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type BestrecipeResult struct {
	ID                      uint
	RecipeName              string
	AmountProduced          uint
	ProductionTime          uint
	MachineName             string
	MachineSpeed            uint
	Rate                    float32
	MachinePowerConsumption uint64
	MachineId               uint
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
	TotalPowerConsumedkW       uint64
	RequiredResourcesPerSecond map[string]float32
	ProducedResourcesPerSecond map[string]float32
	SourceNodes                []int
}

type ResourceSource struct {
	NodeId                          int
	ExcessResourceName              string
	ExcessProducedResourcePerSecond float32
}

type ProductionTree struct {
	TreeNodes                []*ProductionTreeNode
	TargetResource           string
	TargetResourceRate       float32
	TargetResourceSourceNode int
	ExcessResources          []*ResourceSource
}

func Calculate(ctx context.Context, userId int, desiredResourceName string, desiredRate float32, recipes_names []string, machines_names []string, db *sql.DB) ([]byte, error) {
	excessResources := list.New()
	calculationResult := ProductionTree{TreeNodes: make([]*ProductionTreeNode, 0), ExcessResources: make([]*ResourceSource, excessResources.Len())}
	var err error
	calculationResult.TargetResourceSourceNode, err = findAndComputeBestrecipeForResource(ctx, userId, desiredResourceName, desiredRate, recipes_names, machines_names, &calculationResult.TreeNodes, excessResources, db)
	if err != nil {
		return nil, fmt.Errorf("could not compute production chain for resource '%s': %w", desiredResourceName, err)
	}
	for e := excessResources.Front(); e != nil; e = e.Next() {
		newEntry := ResourceSource{NodeId: e.Value.(ResourceSource).NodeId, ExcessResourceName: e.Value.(ResourceSource).ExcessResourceName, ExcessProducedResourcePerSecond: e.Value.(ResourceSource).ExcessProducedResourcePerSecond}
		calculationResult.ExcessResources = append(calculationResult.ExcessResources, &newEntry)
	}
	calculationResult.TargetResource = desiredResourceName
	calculationResult.TargetResourceRate = desiredRate
	byteJSONRepresentation, err := json.Marshal(calculationResult)
	if err != nil {
		return nil, fmt.Errorf("could not generate json representation: %w", err)
	}
	return byteJSONRepresentation, nil
}

func findBestrecipe(ctx context.Context, userId int, desiredResourceName string, recipes_names []string, machines_names []string, db *sql.DB, bestrecipe *BestrecipeResult) error {
	var query string = `SELECT rcp.id, rcp.name AS recipe_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate, m.power_consumption_kw AS machine_power_consumption, m.id AS machine_id
							FROM recipes rcp
							JOIN recipes_outputs ro ON rcp.id = ro.recipes_id
							JOIN resources r ON ro.resources_id = r.id
							JOIN machines_recipes mr ON rcp.id = mr.recipes_id
							JOIN machines m ON mr.machines_id = m.id
							WHERE r.name = '` + desiredResourceName + `'`
	if len(recipes_names) == 0 {
		query += ` AND rcp.default_choice = TRUE `
	} else {
		query += " AND (rcp.default_choice = TRUE OR rcp.name IN ("
		for i, name := range recipes_names {
			if i != 0 {
				query += ","
			}
			query += "'" + name + "'"
		}
		query += ")) "
	}
	if len(machines_names) == 0 {
		query += ` AND m.default_choice = TRUE `
	} else {
		query += " AND (m.default_choice = TRUE OR m.name IN ("
		for i, name := range machines_names {
			if i != 0 {
				query += ","
			}
			query += "'" + name + "'"
		}
		query += ")) "
	}
	query += ` AND rcp.users_id = '` + fmt.Sprint(userId) + `'
			AND ro.users_id = '` + fmt.Sprint(userId) + `'
			AND r.users_id = '` + fmt.Sprint(userId) + `'
			AND mr.users_id = '` + fmt.Sprint(userId) + `'
			AND m.users_id = '` + fmt.Sprint(userId) + `'
			ORDER BY rcp.default_choice, rate DESC
			LIMIT 1;`
	err := db.QueryRowContext(ctx, query).Scan(&(*bestrecipe).ID, &(*bestrecipe).RecipeName, &(*bestrecipe).AmountProduced, &(*bestrecipe).ProductionTime, &(*bestrecipe).MachineName, &(*bestrecipe).MachineSpeed, &(*bestrecipe).Rate, &(*bestrecipe).MachinePowerConsumption, &(*bestrecipe).MachineId)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not find recipe for '%s'", desiredResourceName)
	}
	if err != nil {
		return err
	}
	return nil
}

func getRequiredResources(ctx context.Context, userId int, recipe_id uint, machine_id uint, requiredResources *list.List, db *sql.DB) error {
	var resource RequiredResourceResult
	var query string = `SELECT r.id, r.name AS resource_name, ri.amount AS amount_required, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ri.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipes rcp
			JOIN recipes_inputs ri ON rcp.id = ri.recipes_id
			JOIN resources r ON ri.resources_id = r.id
			JOIN machines_recipes mr ON rcp.id = mr.recipes_id
			JOIN machines m ON mr.machines_id = m.id
			WHERE rcp.id = ?
			AND m.id = ? 
			AND rcp.users_id = '` + fmt.Sprint(userId) + `'
			AND ri.users_id = '` + fmt.Sprint(userId) + `'
			AND r.users_id = '` + fmt.Sprint(userId) + `'
			AND mr.users_id = '` + fmt.Sprint(userId) + `'
			AND m.users_id = '` + fmt.Sprint(userId) + `';`
	rows, err := db.QueryContext(ctx, query, recipe_id, machine_id)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&resource.ID, &resource.ResourceName, &resource.AmountRequired, &resource.ProductionTime, &resource.MachineName, &resource.MachineSpeed, &resource.Rate)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		requiredResources.PushBack(resource)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func getProducedResources(ctx context.Context, userId int, recipe_id uint, machine_id uint, producedResources *list.List, db *sql.DB) error {
	var resource ProducedResourceResult
	var query string = `SELECT r.id, r.name AS resource_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipes rcp
			JOIN recipes_outputs ro ON rcp.id = ro.recipes_id
			JOIN resources r ON ro.resources_id = r.id
			JOIN machines_recipes mr ON rcp.id = mr.recipes_id
			JOIN machines m ON mr.machines_id = m.id
			WHERE rcp.id = ?
			AND m.id = ?
			AND rcp.users_id = '` + fmt.Sprint(userId) + `'
			AND ro.users_id = '` + fmt.Sprint(userId) + `'
			AND r.users_id = '` + fmt.Sprint(userId) + `'
			AND mr.users_id = '` + fmt.Sprint(userId) + `'
			AND m.users_id = '` + fmt.Sprint(userId) + `';`
	rows, err := db.QueryContext(ctx, query, recipe_id, machine_id)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&resource.ID, &resource.ResourceName, &resource.AmountProduced, &resource.ProductionTime, &resource.MachineName, &resource.MachineSpeed, &resource.Rate)
		if err != nil {
			return err
		}
		producedResources.PushBack(resource)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func findAndComputeBestrecipeForResource(ctx context.Context, userId int, desiredResourceName string, desiredRate float32, recipes_names []string, machines_names []string, ProductionTreeNodes *[]*ProductionTreeNode, ExcessResources *list.List, db *sql.DB) (int, error) {
	var bestrecipe BestrecipeResult
	var machinesRequired float32
	var NewNode ProductionTreeNode = ProductionTreeNode{RequiredResourcesPerSecond: make(map[string]float32), ProducedResourcesPerSecond: make(map[string]float32)}
	var RequiredResourcesTemp = make(map[string]float32)
	err := findBestrecipe(ctx, userId, desiredResourceName, recipes_names, machines_names, db, &bestrecipe)
	if err != nil {
		return -1, fmt.Errorf("could not compute production chain for resource '%s': %w", desiredResourceName, err)
	}
	machinesRequired = desiredRate / bestrecipe.Rate
	requiredResources := list.New()
	producedResources := list.New()
	err = getRequiredResources(ctx, userId, bestrecipe.ID, bestrecipe.MachineId, requiredResources, db)
	if err != nil {
		return -1, fmt.Errorf("could not find required resources for recipe '%s': %w", bestrecipe.RecipeName, err)
	}
	err = getProducedResources(ctx, userId, bestrecipe.ID, bestrecipe.MachineId, producedResources, db)
	if err != nil {
		return -1, fmt.Errorf("could not find produced resources for recipe '%s': %w", bestrecipe.RecipeName, err)
	}
	for e := requiredResources.Front(); e != nil; e = e.Next() {
		inserted := false
		NewNode.RequiredResourcesPerSecond[e.Value.(RequiredResourceResult).ResourceName] = e.Value.(RequiredResourceResult).Rate * machinesRequired
		for ei := ExcessResources.Front(); ei != nil; ei = ei.Next() {
			if ei.Value.(ResourceSource).ExcessResourceName == e.Value.(RequiredResourceResult).ResourceName {
				if ei.Value.(ResourceSource).ExcessProducedResourcePerSecond < e.Value.(RequiredResourceResult).Rate*machinesRequired {
					RequiredResourcesTemp[e.Value.(RequiredResourceResult).ResourceName] = (e.Value.(RequiredResourceResult).Rate * machinesRequired) - ei.Value.(ResourceSource).ExcessProducedResourcePerSecond
				} else if ei.Value.(ResourceSource).ExcessProducedResourcePerSecond >= e.Value.(RequiredResourceResult).Rate*machinesRequired {
					newEiElement := ResourceSource{NodeId: ei.Value.(ResourceSource).NodeId, ExcessResourceName: ei.Value.(ResourceSource).ExcessResourceName, ExcessProducedResourcePerSecond: ei.Value.(ResourceSource).ExcessProducedResourcePerSecond - e.Value.(RequiredResourceResult).Rate*machinesRequired}
					if newEiElement.ExcessProducedResourcePerSecond > 0 {
						ExcessResources.InsertBefore(newEiElement, ei)
					}
				}
				ExcessResources.Remove(ei)
				inserted = true
				NewNode.SourceNodes = append(NewNode.SourceNodes, (ei.Value.(ResourceSource).NodeId))
			}
		}

		if !inserted {
			RequiredResourcesTemp[e.Value.(RequiredResourceResult).ResourceName] = e.Value.(RequiredResourceResult).Rate * machinesRequired
		}
	}
	for e := producedResources.Front(); e != nil; e = e.Next() {
		NewNode.ProducedResourcesPerSecond[e.Value.(ProducedResourceResult).ResourceName] = e.Value.(ProducedResourceResult).Rate * machinesRequired
		if e.Value.(ProducedResourceResult).ResourceName != desiredResourceName {
			var excessResource ResourceSource
			excessResource.NodeId = len(*ProductionTreeNodes)
			excessResource.ExcessResourceName = e.Value.(ProducedResourceResult).ResourceName
			excessResource.ExcessProducedResourcePerSecond = e.Value.(ProducedResourceResult).Rate * machinesRequired
			ExcessResources.PushBack(excessResource)
		}
	}
	NewNode.MachineName = bestrecipe.MachineName
	NewNode.MachineNumber = machinesRequired
	NewNode.RecipeName = bestrecipe.RecipeName
	NewNode.NodeId = len(*ProductionTreeNodes)
	NewNode.TotalPowerConsumedkW = uint64(machinesRequired) * bestrecipe.MachinePowerConsumption
	*ProductionTreeNodes = append(*ProductionTreeNodes, &NewNode)
	for resourceName, requiredAmount := range RequiredResourcesTemp {
		if requiredAmount > 0 {
			var sourceNode int
			sourceNode, err = findAndComputeBestrecipeForResource(ctx, userId, resourceName, requiredAmount, recipes_names, machines_names, ProductionTreeNodes, ExcessResources, db)
			if err != nil {
				return -1, fmt.Errorf("could not compute production chain for resource '%s': %w", resourceName, err)
			}
			NewNode.SourceNodes = append(NewNode.SourceNodes, sourceNode)
		}
	}
	return NewNode.NodeId, nil
}
