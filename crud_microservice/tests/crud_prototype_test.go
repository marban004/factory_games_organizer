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

package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/marban004/factory_games_organizer/prototypes"
	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (uts *UnitTestSuite) SetupSuite() {
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = "pH082C./"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "users_data"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		uts.FailNowf("unable to connect to database, error: %s", err.Error())
	}
	uts.db = db
	setupDatabaseSchema(uts)
	cleanDatabaseTables(uts)
}

func (uts *UnitTestSuite) TearDownSuite() {
	teardownDB(uts)
}

func (uts *UnitTestSuite) TearDownTest() {
	cleanDatabaseTables(uts)
}

func (uts *UnitTestSuite) TestSelectMachinesById() {

	expectedRows := []prototypes.MachineInfo{
		{Id: 1, Name: "harvester_mk1", UsersId: 1, InputsSolid: 0, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 20000, DefaultChoice: 1},
		{Id: 4, Name: "assembler_mk1", UsersId: 1, InputsSolid: 2, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 30000, DefaultChoice: 1},
	}
	returnedRows, err := prototypes.SelectMachinesById(context.Background(), uts.db, []int{1, 4}, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectMachines() {

	expectedRows := []prototypes.MachineInfo{
		{Id: 1, Name: "harvester_mk1", UsersId: 1, InputsSolid: 0, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 20000, DefaultChoice: 1},
		{Id: 2, Name: "smelter_mk1", UsersId: 1, InputsSolid: 1, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 10000, DefaultChoice: 1},
		{Id: 3, Name: "constructor_mk1", UsersId: 1, InputsSolid: 1, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 10000, DefaultChoice: 1},
		{Id: 4, Name: "assembler_mk1", UsersId: 1, InputsSolid: 2, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 30000, DefaultChoice: 1},
	}
	returnedRows, err := prototypes.SelectMachines(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestInsertMachines() {
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	input := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := prototypes.InsertMachines(context.Background(), uts.db, input.MachinesList)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectMachines(context.Background(), uts.db, 5, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, input.MachinesList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestUpdateMachines() {
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	update := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := prototypes.UpdateMachines(context.Background(), uts.db, update.MachinesList)
	uts.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		uts.Nil(err)
		rowsChanged += temp
	}

	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectMachines(context.Background(), uts.db, 3, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, update.MachinesList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestDeleteMachines() {
	expectedRows := []prototypes.MachineInfo{
		{Id: 1, Name: "harvester_mk1", UsersId: 1, InputsSolid: 0, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 20000, DefaultChoice: 1},
		{Id: 3, Name: "constructor_mk1", UsersId: 1, InputsSolid: 1, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 10000, DefaultChoice: 1},
	}
	ids := []int{2, 4}
	result, err := prototypes.DeleteMachines(context.Background(), uts.db, ids, 1)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectMachines(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectResourcesById() {

	expectedRows := []prototypes.ResourceInfo{
		{Id: 5, Name: "screw", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 6, Name: "reinforced_iron_plate", UsersId: 1, Liquid: 0, ResourceUnit: ""},
	}
	returnedRows, err := prototypes.SelectResourcesById(context.Background(), uts.db, []int{5, 6}, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectResources() {

	expectedRows := []prototypes.ResourceInfo{
		{Id: 1, Name: "iron_ore", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 2, Name: "iron_ingot", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 3, Name: "iron_plate", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 4, Name: "iron_rod", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 5, Name: "screw", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 6, Name: "reinforced_iron_plate", UsersId: 1, Liquid: 0, ResourceUnit: ""},
	}
	returnedRows, err := prototypes.SelectResources(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestInsertResources() {
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	input := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := prototypes.InsertResources(context.Background(), uts.db, input.ResourcesList)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectResources(context.Background(), uts.db, 7, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, input.ResourcesList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestUpdateResources() {
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	update := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := prototypes.UpdateResources(context.Background(), uts.db, update.ResourcesList)
	uts.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		uts.Nil(err)
		rowsChanged += temp
	}

	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectResources(context.Background(), uts.db, 5, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, update.ResourcesList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestDeleteResources() {
	expectedRows := []prototypes.ResourceInfo{
		{Id: 2, Name: "iron_ingot", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 5, Name: "screw", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 6, Name: "reinforced_iron_plate", UsersId: 1, Liquid: 0, ResourceUnit: ""},
	}
	ids := []int{1, 3, 4}
	result, err := prototypes.DeleteResources(context.Background(), uts.db, ids, 1)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(3), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectResources(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectRecipesById() {

	expectedRows := []prototypes.RecipeInfo{
		{Id: 2, Name: "iron_ingot", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 3, Name: "iron_plate", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 4, Name: "iron_rods", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 5, Name: "screw", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
	}
	returnedRows, err := prototypes.SelectRecipesById(context.Background(), uts.db, []int{2, 3, 4, 5}, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectRecipes() {

	expectedRows := []prototypes.RecipeInfo{
		{Id: 1, Name: "iron_ore_harvesting_default", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 2, Name: "iron_ingot", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 3, Name: "iron_plate", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 4, Name: "iron_rods", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 5, Name: "screw", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 6, Name: "reinforced_iron_plate", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
	}
	returnedRows, err := prototypes.SelectRecipes(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestInsertRecipes() {
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	input := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := prototypes.InsertRecipes(context.Background(), uts.db, input.RecipesList)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectRecipes(context.Background(), uts.db, 7, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, input.RecipesList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestUpdateRecipes() {
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	update := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := prototypes.UpdateRecipes(context.Background(), uts.db, update.RecipesList)
	uts.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		uts.Nil(err)
		rowsChanged += temp
	}

	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectRecipes(context.Background(), uts.db, 5, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, update.RecipesList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestDeleteRecipes() {
	expectedRows := []prototypes.RecipeInfo{
		{Id: 1, Name: "iron_ore_harvesting_default", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 2, Name: "iron_ingot", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 3, Name: "iron_plate", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
	}
	ids := []int{4, 5, 6}
	result, err := prototypes.DeleteRecipes(context.Background(), uts.db, ids, 1)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(3), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectRecipes(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectRecipesInputsById() {

	expectedRows := []prototypes.RecipeInputOutputInfo{
		{Id: 1, UsersId: 1, RecipesId: 2, ResourcesId: 1, Amount: 30},
		{Id: 5, UsersId: 1, RecipesId: 6, ResourcesId: 3, Amount: 30},
	}
	returnedRows, err := prototypes.SelectRecipesInputsById(context.Background(), uts.db, []int{1, 5}, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectRecipesInputs() {

	expectedRows := []prototypes.RecipeInputOutputInfo{
		{Id: 1, UsersId: 1, RecipesId: 2, ResourcesId: 1, Amount: 30},
		{Id: 2, UsersId: 1, RecipesId: 3, ResourcesId: 2, Amount: 30},
		{Id: 3, UsersId: 1, RecipesId: 4, ResourcesId: 2, Amount: 15},
		{Id: 4, UsersId: 1, RecipesId: 5, ResourcesId: 4, Amount: 10},
		{Id: 5, UsersId: 1, RecipesId: 6, ResourcesId: 3, Amount: 30},
		{Id: 6, UsersId: 1, RecipesId: 6, ResourcesId: 5, Amount: 60},
	}
	returnedRows, err := prototypes.SelectRecipesInputs(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestInsertRecipesInputs() {
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	input := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := prototypes.InsertRecipesInputs(context.Background(), uts.db, input.RecipesInputsList)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectRecipesInputs(context.Background(), uts.db, 7, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, input.RecipesInputsList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestUpdateRecipesInputs() {
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	update := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := prototypes.UpdateRecipesInputs(context.Background(), uts.db, update.RecipesInputsList)
	uts.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		uts.Nil(err)
		rowsChanged += temp
	}

	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectRecipesInputs(context.Background(), uts.db, 3, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, update.RecipesInputsList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestDeleteRecipesInputs() {
	expectedRows := []prototypes.RecipeInputOutputInfo{
		{Id: 2, UsersId: 1, RecipesId: 3, ResourcesId: 2, Amount: 30},
		{Id: 3, UsersId: 1, RecipesId: 4, ResourcesId: 2, Amount: 15},
		{Id: 4, UsersId: 1, RecipesId: 5, ResourcesId: 4, Amount: 10},
		{Id: 5, UsersId: 1, RecipesId: 6, ResourcesId: 3, Amount: 30},
	}
	ids := []int{1, 6}
	result, err := prototypes.DeleteRecipesInputs(context.Background(), uts.db, ids, 1)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectRecipesInputs(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectRecipesOutputsById() {

	expectedRows := []prototypes.RecipeInputOutputInfo{
		{Id: 6, UsersId: 1, RecipesId: 6, ResourcesId: 6, Amount: 5},
	}
	returnedRows, err := prototypes.SelectRecipesOutputsById(context.Background(), uts.db, []int{6}, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectRecipesOutputs() {

	expectedRows := []prototypes.RecipeInputOutputInfo{
		{Id: 1, UsersId: 1, RecipesId: 1, ResourcesId: 1, Amount: 60},
		{Id: 2, UsersId: 1, RecipesId: 2, ResourcesId: 2, Amount: 30},
		{Id: 3, UsersId: 1, RecipesId: 3, ResourcesId: 3, Amount: 20},
		{Id: 4, UsersId: 1, RecipesId: 4, ResourcesId: 4, Amount: 15},
		{Id: 5, UsersId: 1, RecipesId: 5, ResourcesId: 5, Amount: 40},
		{Id: 6, UsersId: 1, RecipesId: 6, ResourcesId: 6, Amount: 5},
	}
	returnedRows, err := prototypes.SelectRecipesOutputs(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestInsertRecipesOutputs() {
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	input := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := prototypes.InsertRecipesOutputs(context.Background(), uts.db, input.RecipesOutputsList)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectRecipesOutputs(context.Background(), uts.db, 7, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, input.RecipesOutputsList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestUpdateRecipesOutputs() {
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	update := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := prototypes.UpdateRecipesOutputs(context.Background(), uts.db, update.RecipesOutputsList)
	uts.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		uts.Nil(err)
		rowsChanged += temp
	}

	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectRecipesOutputs(context.Background(), uts.db, 4, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, update.RecipesOutputsList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestDeleteRecipesOutputs() {
	expectedRows := []prototypes.RecipeInputOutputInfo{
		{Id: 1, UsersId: 1, RecipesId: 1, ResourcesId: 1, Amount: 60},
		{Id: 2, UsersId: 1, RecipesId: 2, ResourcesId: 2, Amount: 30},
		{Id: 3, UsersId: 1, RecipesId: 3, ResourcesId: 3, Amount: 20},
		{Id: 4, UsersId: 1, RecipesId: 4, ResourcesId: 4, Amount: 15},
		{Id: 5, UsersId: 1, RecipesId: 5, ResourcesId: 5, Amount: 40},
	}
	ids := []int{6}
	result, err := prototypes.DeleteRecipesOutputs(context.Background(), uts.db, ids, 1)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(1), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectRecipesOutputs(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectMachinesRecipesById() {

	expectedRows := []prototypes.MachinesRecipesInfo{
		{Id: 2, UsersId: 1, RecipesId: 2, MachinesId: 2},
		{Id: 3, UsersId: 1, RecipesId: 3, MachinesId: 3},
		{Id: 5, UsersId: 1, RecipesId: 5, MachinesId: 3},
	}
	returnedRows, err := prototypes.SelectMachinesRecipesById(context.Background(), uts.db, []int{2, 3, 5}, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestSelectMachinesRecipes() {

	expectedRows := []prototypes.MachinesRecipesInfo{
		{Id: 1, UsersId: 1, RecipesId: 1, MachinesId: 1},
		{Id: 2, UsersId: 1, RecipesId: 2, MachinesId: 2},
		{Id: 3, UsersId: 1, RecipesId: 3, MachinesId: 3},
		{Id: 4, UsersId: 1, RecipesId: 4, MachinesId: 3},
		{Id: 5, UsersId: 1, RecipesId: 5, MachinesId: 3},
		{Id: 6, UsersId: 1, RecipesId: 6, MachinesId: 4},
	}
	returnedRows, err := prototypes.SelectMachinesRecipes(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestInsertMachinesRecipes() {
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	input := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := prototypes.InsertMachinesRecipes(context.Background(), uts.db, input.MachinesRecipesList)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectMachinesRecipes(context.Background(), uts.db, 7, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, input.MachinesRecipesList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestUpdateMachinesRecipes() {
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		uts.FailNowf("failed to read file", err.Error())
	}
	update := prototypes.JSONInput{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := prototypes.UpdateMachinesRecipes(context.Background(), uts.db, update.MachinesRecipesList)
	uts.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		uts.Nil(err)
		rowsChanged += temp
	}

	uts.Nil(err)
	uts.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectMachinesRecipes(context.Background(), uts.db, 5, 2, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, update.MachinesRecipesList, "The returned and expected values don't match")
}

func (uts *UnitTestSuite) TestDeleteMachinesRecipes() {
	expectedRows := []prototypes.MachinesRecipesInfo{
		{Id: 1, UsersId: 1, RecipesId: 1, MachinesId: 1},
		{Id: 2, UsersId: 1, RecipesId: 2, MachinesId: 2},
		{Id: 3, UsersId: 1, RecipesId: 3, MachinesId: 3},
		{Id: 4, UsersId: 1, RecipesId: 4, MachinesId: 3},
		{Id: 5, UsersId: 1, RecipesId: 5, MachinesId: 3},
	}
	ids := []int{6}
	result, err := prototypes.DeleteMachinesRecipes(context.Background(), uts.db, ids, 1)
	uts.Nil(err)

	rowsChanged, err := result.RowsAffected()
	uts.Nil(err)
	uts.Equal(int64(1), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := prototypes.SelectMachinesRecipes(context.Background(), uts.db, 0, 0, 1)
	uts.Nil(err)
	uts.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func setupDatabaseSchema(uts *UnitTestSuite) {
	uts.T().Log("setting up database schema")
	_, err := uts.db.Exec(`CREATE DATABASE users_data_test`)
	if err != nil {
		uts.FailNowf("unable to create database", err.Error())
	}

	_, err = uts.db.Exec(`USE users_data_test`)
	if err != nil {
		uts.FailNowf("unable to switch currently used database", err.Error())
	}

	contents, err := os.ReadFile("schema_mysql.sql")
	if err != nil {
		uts.FailNowf("unable to read schema sql file", err.Error())
	}
	commands := strings.Split(string(contents), ";")
	for _, command := range commands {
		if len(command) <= 0 {
			break
		}
		_, err = uts.db.Exec(command)
		if err != nil {
			uts.FailNowf("unable to setup database schema", err.Error())
		}
	}
}

func cleanDatabaseTables(uts *UnitTestSuite) {
	uts.T().Log("cleaning database tables")
	contents, err := os.ReadFile("data_mysql.sql")
	if err != nil {
		uts.FailNowf("unable to read data sql file", err.Error())
	}

	commands := strings.Split(string(contents), ";")
	for _, command := range commands {
		if len(command) <= 0 {
			break
		}
		_, err = uts.db.Exec(command)
		if err != nil {
			uts.FailNowf("unable to setup database contents", err.Error())
		}
	}
}

func teardownDB(uts *UnitTestSuite) {
	_, err := uts.db.Exec("DROP DATABASE users_data_test")
	if err != nil {
		uts.FailNowf("failed to drop test db", err.Error())
	}
	err = uts.db.Close()
	if err != nil {
		uts.FailNowf("failed to close test db", err.Error())
	}
}
