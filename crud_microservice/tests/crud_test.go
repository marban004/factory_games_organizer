package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/marban004/factory_games_organizer/handler"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/model"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/machine"
	machinerecipe "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/machine_recipe"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe"
	recipeinput "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe_input"
	recipeoutput "github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/recipe_output"
	"github.com/marban004/factory_games_organizer/microservice_logic_crud/repository/resource"
	"github.com/stretchr/testify/suite"
)

type CrudIntegrationTestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestCrudIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &CrudIntegrationTestSuite{})
}

func (cits *CrudIntegrationTestSuite) SetupSuite() {
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = "pH082C./"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "users_data"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		cits.FailNowf("unable to connect to database, error: %s", err.Error())
	}
	cits.db = db
	setupDatabaseSchemaCITS(cits)
	cleanDatabaseTablesCITS(cits)
}

func (cits *CrudIntegrationTestSuite) TearDownSuite() {
	teardownDBCITS(cits)
}

func (cits *CrudIntegrationTestSuite) TearDownTest() {
	cleanDatabaseTablesCITS(cits)
}

func (cits *CrudIntegrationTestSuite) TestSelectMachinesById() {
	repo := machine.MySQLRepo{DB: cits.db}
	expectedRows := []model.MachineInfo{
		{Id: 1, Name: "harvester_mk1", UsersId: 1, InputsSolid: 0, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 20000, DefaultChoice: 1},
		{Id: 4, Name: "assembler_mk1", UsersId: 1, InputsSolid: 2, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 30000, DefaultChoice: 1},
	}
	returnedRows, err := repo.SelectMachinesById(context.Background(), []int{1, 4}, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectMachines() {
	repo := machine.MySQLRepo{DB: cits.db}
	expectedRows := []model.MachineInfo{
		{Id: 1, Name: "harvester_mk1", UsersId: 1, InputsSolid: 0, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 20000, DefaultChoice: 1},
		{Id: 2, Name: "smelter_mk1", UsersId: 1, InputsSolid: 1, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 10000, DefaultChoice: 1},
		{Id: 3, Name: "constructor_mk1", UsersId: 1, InputsSolid: 1, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 10000, DefaultChoice: 1},
		{Id: 4, Name: "assembler_mk1", UsersId: 1, InputsSolid: 2, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 30000, DefaultChoice: 1},
	}
	returnedRows, err := repo.SelectMachines(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestInsertMachines() {
	repo := machine.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	input := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := repo.InsertMachines(context.Background(), input.MachinesList)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectMachines(context.Background(), 5, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, input.MachinesList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestUpdateMachines() {
	repo := machine.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	update := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := repo.UpdateMachines(context.Background(), update.MachinesList)
	cits.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		cits.Nil(err)
		rowsChanged += temp
	}

	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectMachines(context.Background(), 3, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, update.MachinesList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestDeleteMachines() {
	repo := machine.MySQLRepo{DB: cits.db}
	expectedRows := []model.MachineInfo{
		{Id: 1, Name: "harvester_mk1", UsersId: 1, InputsSolid: 0, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 20000, DefaultChoice: 1},
		{Id: 3, Name: "constructor_mk1", UsersId: 1, InputsSolid: 1, InputsLiquid: 0, OutputsSolid: 1, OutputsLiquid: 0, Speed: 1, PowerConsumptionKw: 10000, DefaultChoice: 1},
	}
	ids := []int{2, 4}
	result, err := repo.DeleteMachines(context.Background(), ids, 1)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectMachines(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectResourcesById() {
	repo := resource.MySQLRepo{DB: cits.db}
	expectedRows := []model.ResourceInfo{
		{Id: 5, Name: "screw", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 6, Name: "reinforced_iron_plate", UsersId: 1, Liquid: 0, ResourceUnit: ""},
	}
	returnedRows, err := repo.SelectResourcesById(context.Background(), []int{5, 6}, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectResources() {
	repo := resource.MySQLRepo{DB: cits.db}
	expectedRows := []model.ResourceInfo{
		{Id: 1, Name: "iron_ore", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 2, Name: "iron_ingot", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 3, Name: "iron_plate", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 4, Name: "iron_rod", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 5, Name: "screw", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 6, Name: "reinforced_iron_plate", UsersId: 1, Liquid: 0, ResourceUnit: ""},
	}
	returnedRows, err := repo.SelectResources(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestInsertResources() {
	repo := resource.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	input := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := repo.InsertResources(context.Background(), input.ResourcesList)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectResources(context.Background(), 7, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, input.ResourcesList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestUpdateResources() {
	repo := resource.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	update := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := repo.UpdateResources(context.Background(), update.ResourcesList)
	cits.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		cits.Nil(err)
		rowsChanged += temp
	}

	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectResources(context.Background(), 5, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, update.ResourcesList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestDeleteResources() {
	repo := resource.MySQLRepo{DB: cits.db}
	expectedRows := []model.ResourceInfo{
		{Id: 2, Name: "iron_ingot", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 5, Name: "screw", UsersId: 1, Liquid: 0, ResourceUnit: ""},
		{Id: 6, Name: "reinforced_iron_plate", UsersId: 1, Liquid: 0, ResourceUnit: ""},
	}
	ids := []int{1, 3, 4}
	result, err := repo.DeleteResources(context.Background(), ids, 1)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(3), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectResources(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectRecipesById() {
	repo := recipe.MySQLRepo{DB: cits.db}
	expectedRows := []model.RecipeInfo{
		{Id: 2, Name: "iron_ingot", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 3, Name: "iron_plate", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 4, Name: "iron_rods", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 5, Name: "screw", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
	}
	returnedRows, err := repo.SelectRecipesById(context.Background(), []int{2, 3, 4, 5}, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectRecipes() {
	repo := recipe.MySQLRepo{DB: cits.db}
	expectedRows := []model.RecipeInfo{
		{Id: 1, Name: "iron_ore_harvesting_default", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 2, Name: "iron_ingot", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 3, Name: "iron_plate", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 4, Name: "iron_rods", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 5, Name: "screw", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 6, Name: "reinforced_iron_plate", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
	}
	returnedRows, err := repo.SelectRecipes(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestInsertRecipes() {
	repo := recipe.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	input := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := repo.InsertRecipes(context.Background(), input.RecipesList)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectRecipes(context.Background(), 7, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, input.RecipesList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestUpdateRecipes() {
	repo := recipe.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	update := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := repo.UpdateRecipes(context.Background(), update.RecipesList)
	cits.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		cits.Nil(err)
		rowsChanged += temp
	}

	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectRecipes(context.Background(), 5, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, update.RecipesList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestDeleteRecipes() {
	repo := recipe.MySQLRepo{DB: cits.db}
	expectedRows := []model.RecipeInfo{
		{Id: 1, Name: "iron_ore_harvesting_default", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 2, Name: "iron_ingot", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
		{Id: 3, Name: "iron_plate", UsersId: 1, ProductionTimeS: 60, DefaultChoice: 1},
	}
	ids := []int{4, 5, 6}
	result, err := repo.DeleteRecipes(context.Background(), ids, 1)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(3), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectRecipes(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectRecipesInputsById() {
	repo := recipeinput.MySQLRepo{DB: cits.db}
	expectedRows := []model.RecipeInputOutputInfo{
		{Id: 1, UsersId: 1, RecipesId: 2, ResourcesId: 1, Amount: 30},
		{Id: 5, UsersId: 1, RecipesId: 6, ResourcesId: 3, Amount: 30},
	}
	returnedRows, err := repo.SelectRecipesInputsById(context.Background(), []int{1, 5}, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectRecipesInputs() {
	repo := recipeinput.MySQLRepo{DB: cits.db}
	expectedRows := []model.RecipeInputOutputInfo{
		{Id: 1, UsersId: 1, RecipesId: 2, ResourcesId: 1, Amount: 30},
		{Id: 2, UsersId: 1, RecipesId: 3, ResourcesId: 2, Amount: 30},
		{Id: 3, UsersId: 1, RecipesId: 4, ResourcesId: 2, Amount: 15},
		{Id: 4, UsersId: 1, RecipesId: 5, ResourcesId: 4, Amount: 10},
		{Id: 5, UsersId: 1, RecipesId: 6, ResourcesId: 3, Amount: 30},
		{Id: 6, UsersId: 1, RecipesId: 6, ResourcesId: 5, Amount: 60},
	}
	returnedRows, err := repo.SelectRecipesInputs(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestInsertRecipesInputs() {
	repo := recipeinput.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	input := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := repo.InsertRecipesInputs(context.Background(), input.RecipesInputsList)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectRecipesInputs(context.Background(), 7, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, input.RecipesInputsList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestUpdateRecipesInputs() {
	repo := recipeinput.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	update := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := repo.UpdateRecipesInputs(context.Background(), update.RecipesInputsList)
	cits.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		cits.Nil(err)
		rowsChanged += temp
	}

	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectRecipesInputs(context.Background(), 3, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, update.RecipesInputsList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestDeleteRecipesInputs() {
	repo := recipeinput.MySQLRepo{DB: cits.db}
	expectedRows := []model.RecipeInputOutputInfo{
		{Id: 2, UsersId: 1, RecipesId: 3, ResourcesId: 2, Amount: 30},
		{Id: 3, UsersId: 1, RecipesId: 4, ResourcesId: 2, Amount: 15},
		{Id: 4, UsersId: 1, RecipesId: 5, ResourcesId: 4, Amount: 10},
		{Id: 5, UsersId: 1, RecipesId: 6, ResourcesId: 3, Amount: 30},
	}
	ids := []int{1, 6}
	result, err := repo.DeleteRecipesInputs(context.Background(), ids, 1)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectRecipesInputs(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectRecipesOutputsById() {
	repo := recipeoutput.MySQLRepo{DB: cits.db}
	expectedRows := []model.RecipeInputOutputInfo{
		{Id: 6, UsersId: 1, RecipesId: 6, ResourcesId: 6, Amount: 5},
	}
	returnedRows, err := repo.SelectRecipesOutputsById(context.Background(), []int{6}, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectRecipesOutputs() {
	repo := recipeoutput.MySQLRepo{DB: cits.db}
	expectedRows := []model.RecipeInputOutputInfo{
		{Id: 1, UsersId: 1, RecipesId: 1, ResourcesId: 1, Amount: 60},
		{Id: 2, UsersId: 1, RecipesId: 2, ResourcesId: 2, Amount: 30},
		{Id: 3, UsersId: 1, RecipesId: 3, ResourcesId: 3, Amount: 20},
		{Id: 4, UsersId: 1, RecipesId: 4, ResourcesId: 4, Amount: 15},
		{Id: 5, UsersId: 1, RecipesId: 5, ResourcesId: 5, Amount: 40},
		{Id: 6, UsersId: 1, RecipesId: 6, ResourcesId: 6, Amount: 5},
	}
	returnedRows, err := repo.SelectRecipesOutputs(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestInsertRecipesOutputs() {
	repo := recipeoutput.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	input := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := repo.InsertRecipesOutputs(context.Background(), input.RecipesOutputsList)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectRecipesOutputs(context.Background(), 7, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, input.RecipesOutputsList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestUpdateRecipesOutputs() {
	repo := recipeoutput.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	update := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := repo.UpdateRecipesOutputs(context.Background(), update.RecipesOutputsList)
	cits.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		cits.Nil(err)
		rowsChanged += temp
	}

	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectRecipesOutputs(context.Background(), 4, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, update.RecipesOutputsList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestDeleteRecipesOutputs() {
	repo := recipeoutput.MySQLRepo{DB: cits.db}
	expectedRows := []model.RecipeInputOutputInfo{
		{Id: 1, UsersId: 1, RecipesId: 1, ResourcesId: 1, Amount: 60},
		{Id: 2, UsersId: 1, RecipesId: 2, ResourcesId: 2, Amount: 30},
		{Id: 3, UsersId: 1, RecipesId: 3, ResourcesId: 3, Amount: 20},
		{Id: 4, UsersId: 1, RecipesId: 4, ResourcesId: 4, Amount: 15},
		{Id: 5, UsersId: 1, RecipesId: 5, ResourcesId: 5, Amount: 40},
	}
	ids := []int{6}
	result, err := repo.DeleteRecipesOutputs(context.Background(), ids, 1)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(1), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectRecipesOutputs(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectMachinesRecipesById() {
	repo := machinerecipe.MySQLRepo{DB: cits.db}
	expectedRows := []model.MachinesRecipesInfo{
		{Id: 2, UsersId: 1, RecipesId: 2, MachinesId: 2},
		{Id: 3, UsersId: 1, RecipesId: 3, MachinesId: 3},
		{Id: 5, UsersId: 1, RecipesId: 5, MachinesId: 3},
	}
	returnedRows, err := repo.SelectMachinesRecipesById(context.Background(), []int{2, 3, 5}, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestSelectMachinesRecipes() {
	repo := machinerecipe.MySQLRepo{DB: cits.db}
	expectedRows := []model.MachinesRecipesInfo{
		{Id: 1, UsersId: 1, RecipesId: 1, MachinesId: 1},
		{Id: 2, UsersId: 1, RecipesId: 2, MachinesId: 2},
		{Id: 3, UsersId: 1, RecipesId: 3, MachinesId: 3},
		{Id: 4, UsersId: 1, RecipesId: 4, MachinesId: 3},
		{Id: 5, UsersId: 1, RecipesId: 5, MachinesId: 3},
		{Id: 6, UsersId: 1, RecipesId: 6, MachinesId: 4},
	}
	returnedRows, err := repo.SelectMachinesRecipes(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestInsertMachinesRecipes() {
	repo := machinerecipe.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_input.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	input := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &input)

	result, err := repo.InsertMachinesRecipes(context.Background(), input.MachinesRecipesList)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectMachinesRecipes(context.Background(), 7, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, input.MachinesRecipesList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestUpdateMachinesRecipes() {
	repo := machinerecipe.MySQLRepo{DB: cits.db}
	jsonFileBytes, err := os.ReadFile("test_update.json")
	if err != nil {
		cits.FailNowf("failed to read file", err.Error())
	}
	update := handler.JSONData{}
	json.Unmarshal(jsonFileBytes, &update)

	resultArr, err := repo.UpdateMachinesRecipes(context.Background(), update.MachinesRecipesList)
	cits.Nil(err)

	rowsChanged := int64(0)
	for _, result := range resultArr {
		temp, err := result.RowsAffected()
		cits.Nil(err)
		rowsChanged += temp
	}

	cits.Nil(err)
	cits.Equal(int64(2), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectMachinesRecipes(context.Background(), 5, 2, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, update.MachinesRecipesList, "The returned and expected values don't match")
}

func (cits *CrudIntegrationTestSuite) TestDeleteMachinesRecipes() {
	repo := machinerecipe.MySQLRepo{DB: cits.db}
	expectedRows := []model.MachinesRecipesInfo{
		{Id: 1, UsersId: 1, RecipesId: 1, MachinesId: 1},
		{Id: 2, UsersId: 1, RecipesId: 2, MachinesId: 2},
		{Id: 3, UsersId: 1, RecipesId: 3, MachinesId: 3},
		{Id: 4, UsersId: 1, RecipesId: 4, MachinesId: 3},
		{Id: 5, UsersId: 1, RecipesId: 5, MachinesId: 3},
	}
	ids := []int{6}
	result, err := repo.DeleteMachinesRecipes(context.Background(), ids, 1)
	cits.Nil(err)

	rowsChanged, err := result.RowsAffected()
	cits.Nil(err)
	cits.Equal(int64(1), rowsChanged, "The number of changed rows differs from expected")

	returnedRows, err := repo.SelectMachinesRecipes(context.Background(), 0, 0, 1)
	cits.Nil(err)
	cits.ElementsMatch(returnedRows, expectedRows, "The returned and expected values don't match")
}

func setupDatabaseSchemaCITS(cits *CrudIntegrationTestSuite) {
	cits.T().Log("setting up database schema")
	_, err := cits.db.Exec(`CREATE DATABASE users_data_test`)
	if err != nil {
		cits.FailNowf("unable to create database", err.Error())
	}

	_, err = cits.db.Exec(`USE users_data_test`)
	if err != nil {
		cits.FailNowf("unable to switch currently used database", err.Error())
	}

	contents, err := os.ReadFile("schema_mysql.sql")
	if err != nil {
		cits.FailNowf("unable to read schema sql file", err.Error())
	}
	commands := strings.Split(string(contents), ";")
	for _, command := range commands {
		if len(command) <= 0 {
			break
		}
		_, err = cits.db.Exec(command)
		if err != nil {
			cits.FailNowf("unable to setup database schema", err.Error())
		}
	}
}

func cleanDatabaseTablesCITS(cits *CrudIntegrationTestSuite) {
	cits.T().Log("cleaning database tables")
	contents, err := os.ReadFile("data_mysql.sql")
	if err != nil {
		cits.FailNowf("unable to read data sql file", err.Error())
	}

	commands := strings.Split(string(contents), ";")
	for _, command := range commands {
		if len(command) <= 0 {
			break
		}
		_, err = cits.db.Exec(command)
		if err != nil {
			cits.FailNowf("unable to setup database contents", err.Error())
		}
	}
}

func teardownDBCITS(cits *CrudIntegrationTestSuite) {
	_, err := cits.db.Exec("DROP DATABASE users_data_test")
	if err != nil {
		cits.FailNowf("failed to drop test db", err.Error())
	}
	err = cits.db.Close()
	if err != nil {
		cits.FailNowf("failed to close test db", err.Error())
	}
}
