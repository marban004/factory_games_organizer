package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/marban004/factory_games_organizer/prototypes"
)

func main() {
	cfg := mysql.NewConfig()
	cfg.User = "crud_microservice"
	cfg.Passwd = "juG56#ian>LK90"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "users_data"
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	jsonFileBytes, err := os.ReadFile("prototypes/test_input.json")
	if err != nil {
		panic(err)
	}
	data := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &data)
	fmt.Printf("%+v", data)
	result, err := prototypes.InsertMachines(context.Background(), db, data.MachinesList)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())

	jsonFileBytes, err = os.ReadFile("prototypes/test_update.json")
	if err != nil {
		panic(err)
	}
	data = prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &data)
	resultList, err := prototypes.UpdateMachines(context.Background(), db, data.MachinesList)
	if err != nil {
		panic(err)
	}
	for _, result := range resultList {
		fmt.Println(result)
	}
	ids := []int{5, 6}
	result, err = prototypes.DeleteMachines(context.Background(), db, ids)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected()) //manually use "ALTER TABLE machines AUTO_INCREMENT = 4;" User crud_microservice cannot modify table structures (no ALTER privilege)
	db.Close()
}
