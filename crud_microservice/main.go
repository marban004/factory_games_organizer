package main

import (
	"bufio"
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
	userId := 1
	if err != nil {
		panic(err)
	}
	jsonFileBytes, err := os.ReadFile("prototypes/test_input.json")
	if err != nil {
		panic(err)
	}
	input := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &input)
	fmt.Printf("%+v\n", input)
	result, err := prototypes.InsertMachines(context.Background(), db, input.MachinesList)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())

	result, err = prototypes.InsertResources(context.Background(), db, input.ResourcesList)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())

	bufio.NewReader(os.Stdin).ReadBytes('\n')

	jsonFileBytes, err = os.ReadFile("prototypes/test_update.json")
	if err != nil {
		panic(err)
	}
	update := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &update)
	resultList, err := prototypes.UpdateMachines(context.Background(), db, update.MachinesList)
	if err != nil {
		panic(err)
	}
	for _, result := range resultList {
		fmt.Println(result.RowsAffected())
	}
	resultList, err = prototypes.UpdateResources(context.Background(), db, update.ResourcesList)
	if err != nil {
		panic(err)
	}
	for _, result := range resultList {
		fmt.Println(result.RowsAffected())
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')

	ids := []int{5, 6}
	result, err = prototypes.DeleteMachines(context.Background(), db, ids, userId)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected()) //relod data by executing schema_mysql.sql, then data_mysql.sql to reset all auto increment sequences

	ids = []int{7, 8}
	result, err = prototypes.DeleteResources(context.Background(), db, ids, userId)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
	db.Close()
}

// ctx context.Context, db *sql.DB, input prototypes.JSONInput, update prototypes.JSONInput, ids []int
