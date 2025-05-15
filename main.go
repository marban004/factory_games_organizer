package main

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/marban004/factory_games_organizer.git/prototypes"
	"gorm.io/gorm"
)

func main() {
	var dbLocation = "F:\\Uczelnia\\PAW\\golang\\factory_games_organizer\\test_db\\test_data.db"
	db, err := gorm.Open(sqlite.Open(dbLocation), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	var altRecipies = [0]string{}
	prototypes.Calculate("reinforced_iron_plate", 0.5, altRecipies[:], db)
	fmt.Println("done")
}
