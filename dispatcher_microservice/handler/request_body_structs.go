package handler

type JSONDataUsers struct {
	UserLogin    string
	UserPassword string
}

type JSONDataCrud struct {
	MachinesList        []MachineInfo
	ResourcesList       []ResourceInfo
	RecipesList         []RecipeInfo
	RecipesInputsList   []RecipeInputOutputInfo
	RecipesOutputsList  []RecipeInputOutputInfo
	MachinesRecipesList []MachinesRecipesInfo
}

type MachineInfo struct {
	Id                 uint
	Name               string
	UsersId            uint
	InputsSolid        uint
	InputsLiquid       uint
	OutputsSolid       uint
	OutputsLiquid      uint
	Speed              float32
	PowerConsumptionKw uint
	DefaultChoice      uint8
}

type MachinesRecipesInfo struct {
	Id         uint
	UsersId    uint
	RecipesId  uint
	MachinesId uint
}

type RecipeInfo struct {
	Id              uint
	Name            string
	UsersId         uint
	ProductionTimeS uint
	DefaultChoice   uint8
}

type RecipeInputOutputInfo struct {
	Id          uint
	UsersId     uint
	RecipesId   uint
	ResourcesId uint
	Amount      uint
}

type ResourceInfo struct {
	Id           uint
	Name         string
	UsersId      uint
	Liquid       uint8
	ResourceUnit string
}

type DeleteInputCrud struct {
	MachinesIds        []int
	ResourcesIds       []int
	RecipesIds         []int
	RecipesInputsIds   []int
	RecipesOutputsIds  []int
	MachinesRecipesIds []int
}
