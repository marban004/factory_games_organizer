package main

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/marban004/factory_games_organizer.git/prototypes"
	"gorm.io/gorm"
)

func main() {
	var dbLocation = "F:\\Uczelnia\\PAW\\golang\\factory_games_organizer\\test_db\\test_data.db"
	var desiredResourceName = "reinforced_iron_plate"
	var userId = 1
	db, err := gorm.Open(sqlite.Open(dbLocation), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	var altRecipies = [0]string{}
	byteJSONRepresentation, err := prototypes.Calculate(userId, desiredResourceName, 0.5, altRecipies[:], db)
	if err != nil {
		fmt.Printf("Could not generate production tree for '%s', reason: %v \n", desiredResourceName, err)
	}
	fmt.Println(string(byteJSONRepresentation))
	fmt.Println("done")
}
