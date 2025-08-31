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
	MicroserviceURl    string
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
