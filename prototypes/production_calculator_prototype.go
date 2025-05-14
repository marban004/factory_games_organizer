package prototypes

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Machine struct {
	id             uint
	name           string
	users_id       uint
	inputs_solid   uint
	inputs_liquid  uint
	outputs_solid  uint
	outputs_liquid uint
	speed          uint
	default_choice uint
}

type Recipies_inputs struct {
	id           uint
	users_id     uint
	recipies_id  uint
	resources_id uint
	amount       uint
}

type Recipies_outputs struct {
	id           uint
	users_id     uint
	recipies_id  uint
	resources_id uint
	amount       uint
}

func Calculate() {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("calculated")
}
