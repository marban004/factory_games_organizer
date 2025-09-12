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

package handler

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type CreateUserResponse struct {
	UsersCreated uint
}

type UpdateUserResponse struct {
	UsersUpdated uint
}

type LoginResponse struct {
	Jwt string
}

type DeleteUserResponse struct {
	UsersDeleted uint
}

type MicroserviceHealth struct {
	MicroserviceURL    string
	MicroserviceStatus string
	DatabaseStatus     string
}

type HealthResponse struct {
	DispatcherStatus       string
	UsersMicroservice      []MicroserviceHealth
	CrudMicroservice       []MicroserviceHealth
	CalculatorMicroservice []MicroserviceHealth
}

type InsertResponseCrud struct {
	MachinesInserted        uint
	ResourcesInserted       uint
	RecipesInserted         uint
	RecipesInputsInserted   uint
	RecipesOutputsInserted  uint
	MachinesRecipesInserted uint
}

type UpdateResponseCrud struct {
	MachinesUpdated        uint
	ResourcesUpdated       uint
	RecipesUpdated         uint
	RecipesInputsUpdated   uint
	RecipesOutputsUpdated  uint
	MachinesRecipesUpdated uint
}

type DeleteResponseCrud struct {
	MachinesDeleted        uint
	ResourcesDeleted       uint
	RecipesDeleted         uint
	RecipesInputsDeleted   uint
	RecipesOutputsDeleted  uint
	MachinesRecipesDeleted uint
}

type ProductionTreeCalculator struct {
	TreeNodes                []*ProductionTreeNode
	TargetResource           string
	TargetResourceRate       float32
	TargetResourceSourceNode int
	ExcessResources          []*ResourceSource
}

type ResourceSource struct {
	NodeId                          int
	ExcessResourceName              string
	ExcessProducedResourcePerSecond float32
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

type StatsResponse struct {
	ApiUsageStats    *orderedmap.OrderedMap[string, map[string]int]
	TrackingPeriodMs int64
	NoPeriods        uint64
}
