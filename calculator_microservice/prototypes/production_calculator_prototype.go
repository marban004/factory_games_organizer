//     This is Factory Games Organizer api. Api is responsible for creating, updating and authenicating api users, CRUD operations on database associated with the api and provides production calculator service.
//     Copyright (C) 2025  Marek Banaś

//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.

//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see https://www.gnu.org/licenses/.

package prototypes

import (
	"container/list"
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

// If using default recipes pass empty string array
func Calculate(userId int, desiredResourceName string, desiredRate float32, names []string, db *sql.DB) ([]byte, error) {
	excessResources := list.New()
	calculationResult := ProductionTree{TreeNodes: make([]*ProductionTreeNode, 0), ExcessResources: make([]*ResourceSource, excessResources.Len())}
	var err error
	calculationResult.TargetResourceSourceNode, err = findAndComputeBestrecipeForResource(userId, desiredResourceName, desiredRate, names, &calculationResult.TreeNodes, excessResources, db)
	if err != nil {
		return nil, fmt.Errorf("could not compute production chain for resource '%s': %w", desiredResourceName, err)
	}
	// for _, Node := range productionTreeNodes {
	// 	fmt.Println("Making recipe:", Node.RecipeName, Node.MachineName, Node.MachineNumber)
	// 	fmt.Println("Node number:", Node.NodeId, "Source nodes: ")
	// 	fmt.Println("------------------------------------------")
	// }
	for e := excessResources.Front(); e != nil; e = e.Next() {
		newEntry := ResourceSource{NodeId: e.Value.(ResourceSource).NodeId, ExcessResourceName: e.Value.(ResourceSource).ExcessResourceName, ExcessProducedResourcePerSecond: e.Value.(ResourceSource).ExcessProducedResourcePerSecond}
		calculationResult.ExcessResources = append(calculationResult.ExcessResources, &newEntry)
	}
	// for _, e := range productionTreeNodes {
	// 	calculationResult.TreeNodes = append(calculationResult.TreeNodes, e)
	// }	calculationResult.TargetResource = desiredResourceName
	calculationResult.TargetResourceRate = desiredRate
	byteJSONRepresentation, err := json.Marshal(calculationResult)
	if err != nil {
		return nil, fmt.Errorf("could not generate json representation: %w", err)
	}
	return byteJSONRepresentation, nil
}

func findBestrecipe(userId int, desiredResourceName string, names []string, db *sql.DB, bestrecipe *BestrecipeResult) error {
	var query string = `SELECT rcp.id, rcp.name AS recipe_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate, m.power_consumption_kw AS machine_power_consumption
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
	query += ` AND rcp.users_id = '` + fmt.Sprint(userId) + `'
			AND ro.users_id = '` + fmt.Sprint(userId) + `'
			AND r.users_id = '` + fmt.Sprint(userId) + `'
			AND mr.users_id = '` + fmt.Sprint(userId) + `'
			AND m.users_id = '` + fmt.Sprint(userId) + `'
			ORDER BY rcp.default_choice, rate DESC
			LIMIT 1;`
	err := db.QueryRow(query).Scan(&(*bestrecipe).ID, &(*bestrecipe).RecipeName, &(*bestrecipe).AmountProduced, &(*bestrecipe).ProductionTime, &(*bestrecipe).MachineName, &(*bestrecipe).MachineSpeed, &(*bestrecipe).Rate, &(*bestrecipe).MachinePowerConsumption)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not find recipe for '%s'", desiredResourceName)
	}
	if err != nil {
		return err
	}
	return nil
	//fmt.Println("selected recipe", bestrecipe.RecipeName)
}

func getRequiredResources(userId int, id uint, requiredResources *list.List, db *sql.DB) error {
	var resource RequiredResourceResult
	var query string = `SELECT r.id, r.name AS resource_name, ri.amount AS amount_required, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ri.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipes rcp
			JOIN recipes_inputs ri ON rcp.id = ri.recipes_id
			JOIN resources r ON ri.resources_id = r.id
			JOIN machines_recipes mr ON rcp.id = mr.recipes_id
			JOIN machines m ON mr.machines_id = m.id
			WHERE rcp.id = ? 
			AND rcp.users_id = '` + fmt.Sprint(userId) + `'
			AND ri.users_id = '` + fmt.Sprint(userId) + `'
			AND r.users_id = '` + fmt.Sprint(userId) + `'
			AND mr.users_id = '` + fmt.Sprint(userId) + `'
			AND m.users_id = '` + fmt.Sprint(userId) + `';`
	rows, err := db.Query(query, id)
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

func getProducedResources(userId int, id uint, producedResources *list.List, db *sql.DB) error {
	var resource ProducedResourceResult
	var query string = `SELECT r.id, r.name AS resource_name, ro.amount AS amount_produced, rcp.production_time_s AS production_time, m.name AS machine_name, m.speed as machine_speed, (CAST(ro.amount AS FLOAT)/rcp.production_time_s*m.speed) rate FROM recipes rcp
			JOIN recipes_outputs ro ON rcp.id = ro.recipes_id
			JOIN resources r ON ro.resources_id = r.id
			JOIN machines_recipes mr ON rcp.id = mr.recipes_id
			JOIN machines m ON mr.machines_id = m.id
			WHERE rcp.id = ?
			AND rcp.users_id = '` + fmt.Sprint(userId) + `'
			AND ro.users_id = '` + fmt.Sprint(userId) + `'
			AND r.users_id = '` + fmt.Sprint(userId) + `'
			AND mr.users_id = '` + fmt.Sprint(userId) + `'
			AND m.users_id = '` + fmt.Sprint(userId) + `';`
	rows, err := db.Query(query, id)
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

func findAndComputeBestrecipeForResource(userId int, desiredResourceName string, desiredRate float32, names []string, ProductionTreeNodes *[]*ProductionTreeNode, ExcessResources *list.List, db *sql.DB) (int, error) {
	//Finding best recipe
	//fmt.Println(desiredResourceName)
	var bestrecipe BestrecipeResult
	var machinesRequired float32
	var NewNode ProductionTreeNode = ProductionTreeNode{RequiredResourcesPerSecond: make(map[string]float32), ProducedResourcesPerSecond: make(map[string]float32)}
	var RequiredResourcesTemp = make(map[string]float32)
	err := findBestrecipe(userId, desiredResourceName, names, db, &bestrecipe)
	if err != nil {
		return -1, fmt.Errorf("could not compute production chain for resource '%s': %w", desiredResourceName, err)
	}
	machinesRequired = desiredRate / bestrecipe.Rate
	//Found best recipe for resource
	//Finding required resources for best recipe
	requiredResources := list.New()
	producedResources := list.New()
	err = getRequiredResources(userId, bestrecipe.ID, requiredResources, db)
	if err != nil {
		return -1, fmt.Errorf("could not find required resources for recipe '%s': %w", bestrecipe.RecipeName, err)
	}
	err = getProducedResources(userId, bestrecipe.ID, producedResources, db)
	if err != nil {
		return -1, fmt.Errorf("could not find produced resources for recipe '%s': %w", bestrecipe.RecipeName, err)
	}
	//fmt.Println(ExcessResources.Front())
	for e := requiredResources.Front(); e != nil; e = e.Next() {
		//fmt.Println(requiredResources.Len(), ExcessResources.Front())
		inserted := false
		//fmt.Println(e.Value.(RequiredResourceResult).ResourceName)
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
		//fmt.Println("produced resources list:", producedResources, e.Value.(ProducedResourceResult).ResourceName)
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
	//fmt.Println(NewNode.recipeName)
	for resourceName, requiredAmount := range RequiredResourcesTemp {
		if requiredAmount > 0 {
			var sourceNode int
			sourceNode, err = findAndComputeBestrecipeForResource(userId, resourceName, requiredAmount, names, ProductionTreeNodes, ExcessResources, db)
			if err != nil {
				return -1, fmt.Errorf("could not compute production chain for resource '%s': %w", resourceName, err)
			}
			NewNode.SourceNodes = append(NewNode.SourceNodes, sourceNode)
		}
	}
	return NewNode.NodeId, nil
}
