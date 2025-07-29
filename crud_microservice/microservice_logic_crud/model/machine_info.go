package model

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
