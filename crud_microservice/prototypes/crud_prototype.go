package prototypes

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

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

type ResourceInfo struct {
	Id           uint
	Name         string
	UsersId      uint
	Liquid       uint8
	ResourceUnit string
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

type JSONInput struct {
	MachinesList       []MachineInfo
	ResourcesList      []ResourceInfo
	RecipesList        []RecipeInfo
	RecipesInputsList  []RecipeInputOutputInfo
	RecipesOutputsList []RecipeInputOutputInfo
}

func SelectMachines(ctx context.Context, db *sql.DB, startId int, rowsRet int) ([]MachineInfo, error) {
	query := "SELECT * FROM machines WHERE id >= " + strconv.Itoa(startId)
	if rowsRet > 0 {
		query += " LIMIT " + strconv.Itoa(rowsRet)
	}
	query += ";"
	result, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from db: %w", err)
	}
	var resultRows []MachineInfo
	for result.Next() {
		var row MachineInfo
		err = result.Scan(&row.Id, &row.Name, &row.UsersId, &row.InputsSolid, &row.InputsLiquid, &row.OutputsSolid, &row.OutputsLiquid, &row.Speed, &row.PowerConsumptionKw, &row.DefaultChoice)
		if err != nil {
			return nil, fmt.Errorf("could not parse data retrieved from db: %w", err)
		}
		resultRows = append(resultRows, row)
	}
	err = result.Err()
	if err != nil {
		return nil, fmt.Errorf("encountered an unexpected error: %w", err)
	}
	return resultRows, nil
}

func InsertMachines(ctx context.Context, db *sql.DB, data []MachineInfo) (sql.Result, error) {
	query := "INSERT INTO machines(name, users_id, inputs_solid, inputs_liquid, outputs_solid, outputs_liquid, speed, power_consumption_kw, default_choice) VALUES"
	for i, entry := range data {
		if i != 0 {
			query += ","
		}
		query += ` ("` + entry.Name +
			`", ` + fmt.Sprint(entry.UsersId) +
			`, ` + fmt.Sprint(entry.InputsSolid) +
			`, ` + fmt.Sprint(entry.InputsLiquid) +
			`, ` + fmt.Sprint(entry.OutputsSolid) +
			`, ` + fmt.Sprint(entry.OutputsLiquid) +
			`, ` + fmt.Sprint(entry.Speed) +
			`, ` + fmt.Sprint(entry.PowerConsumptionKw) +
			`, ` + fmt.Sprint(entry.DefaultChoice) + `)`
	}
	query += ";"
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been inserted: %w", err)
	}
	return result, nil
}

func DeleteMachines(ctx context.Context, db *sql.DB, ids []int) (sql.Result, error) {
	query := "DELETE FROM machines WHERE id in ("
	for i, id := range ids {
		if i != 0 {
			query += ","
		}
		query += " " + fmt.Sprint(id)
	}
	query += ");"
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been deleted: %w", err)
	}
	return result, nil
}

func UpdateMachines(ctx context.Context, db *sql.DB, data []MachineInfo) ([]sql.Result, error) {
	results := []sql.Result{}
	for _, entry := range data {
		query := fmt.Sprintf("UPDATE machines SET name='%s', inputs_solid=%d, inputs_liquid=%d, outputs_solid=%d, outputs_liquid=%d, speed=%f, power_consumption_kw=%d, default_choice=%d WHERE id=%d;",
			entry.Name, entry.InputsSolid, entry.InputsLiquid, entry.OutputsSolid, entry.OutputsLiquid, entry.Speed, entry.PowerConsumptionKw, entry.DefaultChoice, entry.Id)
		result, err := db.ExecContext(ctx, query)
		results = append(results, result)
		if err != nil {
			return results, fmt.Errorf("data has not been fully updated: %w", err)
		}
	}
	return results, nil
}

func SelectResources(ctx context.Context, db *sql.DB, startId int, rowsRet int) ([]ResourceInfo, error) {
	query := "SELECT * FROM resources WHERE id >= " + strconv.Itoa(startId)
	if rowsRet > 0 {
		query += " LIMIT " + strconv.Itoa(rowsRet)
	}
	query += ";"
	result, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from db: %w", err)
	}
	var resultRows []ResourceInfo
	for result.Next() {
		var row ResourceInfo
		err = result.Scan(&row.Id, &row.Name, &row.UsersId, &row.Liquid, &row.ResourceUnit)
		if err != nil {
			return nil, fmt.Errorf("could not parse data retrieved from db: %w", err)
		}
		resultRows = append(resultRows, row)
	}
	err = result.Err()
	if err != nil {
		return nil, fmt.Errorf("encountered an unexpected error: %w", err)
	}
	return resultRows, nil
}

func InsertResources(ctx context.Context, db *sql.DB, data []ResourceInfo) (sql.Result, error) {
	query := "INSERT INTO resources(name, users_id, liquid, resource_unit) VALUES"
	for i, entry := range data {
		if i != 0 {
			query += ","
		}
		query += ` ("` + entry.Name +
			`", ` + fmt.Sprint(entry.UsersId) +
			`, ` + fmt.Sprint(entry.Liquid) +
			`, "` + fmt.Sprint(entry.ResourceUnit) + `")`
	}
	query += ";"
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been inserted: %w", err)
	}
	return result, nil
}

func DeleteResources(ctx context.Context, db *sql.DB, ids []int) (sql.Result, error) {
	query := "DELETE FROM resources WHERE id in ("
	for i, id := range ids {
		if i != 0 {
			query += ","
		}
		query += " " + fmt.Sprint(id)
	}
	query += ");"
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been deleted: %w", err)
	}
	return result, nil
}

func UpdateResources(ctx context.Context, db *sql.DB, data []ResourceInfo) ([]sql.Result, error) {
	results := []sql.Result{}
	for _, entry := range data {
		query := fmt.Sprintf("UPDATE resources SET name='%s', liquid=%d, resource_unit='%s' WHERE id=%d;",
			entry.Name, entry.Liquid, entry.ResourceUnit, entry.Id)
		result, err := db.ExecContext(ctx, query)
		results = append(results, result)
		if err != nil {
			return results, fmt.Errorf("data has not been fully updated: %w", err)
		}
	}
	return results, nil
}

func SelectRecipes(ctx context.Context, db *sql.DB, startId int, rowsRet int) ([]RecipeInfo, error) {
	query := "SELECT * FROM Recipes WHERE id >= " + strconv.Itoa(startId)
	if rowsRet > 0 {
		query += " LIMIT " + strconv.Itoa(rowsRet)
	}
	query += ";"
	result, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from db: %w", err)
	}
	var resultRows []RecipeInfo
	for result.Next() {
		var row RecipeInfo
		err = result.Scan(&row.Id, &row.Name, &row.UsersId, &row.ProductionTimeS, &row.DefaultChoice)
		if err != nil {
			return nil, fmt.Errorf("could not parse data retrieved from db: %w", err)
		}
		resultRows = append(resultRows, row)
	}
	err = result.Err()
	if err != nil {
		return nil, fmt.Errorf("encountered an unexpected error: %w", err)
	}
	return resultRows, nil
}

func InsertRecipes(ctx context.Context, db *sql.DB, data []RecipeInfo) (sql.Result, error) {
	query := "INSERT INTO recipes(name, users_id, production_time_s, default_choice) VALUES"
	for i, entry := range data {
		if i != 0 {
			query += ","
		}
		query += ` ("` + entry.Name +
			`", ` + fmt.Sprint(entry.UsersId) +
			`, ` + fmt.Sprint(entry.ProductionTimeS) +
			`, "` + fmt.Sprint(entry.DefaultChoice) + `")`
	}
	query += ";"
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been inserted: %w", err)
	}
	return result, nil
}

func DeleteRecipes(ctx context.Context, db *sql.DB, ids []int) (sql.Result, error) {
	query := "DELETE FROM recipes WHERE id in ("
	for i, id := range ids {
		if i != 0 {
			query += ","
		}
		query += " " + fmt.Sprint(id)
	}
	query += ");"
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been deleted: %w", err)
	}
	return result, nil
}

func UpdateRecipes(ctx context.Context, db *sql.DB, data []RecipeInfo) ([]sql.Result, error) {
	results := []sql.Result{}
	for _, entry := range data {
		query := fmt.Sprintf("UPDATE recipes SET name='%s', production_time_s=%d, default_choice='%d' WHERE id=%d;",
			entry.Name, entry.ProductionTimeS, entry.DefaultChoice, entry.Id)
		result, err := db.ExecContext(ctx, query)
		results = append(results, result)
		if err != nil {
			return results, fmt.Errorf("data has not been fully updated: %w", err)
		}
	}
	return results, nil
}

func SelectRecipesInputs(ctx context.Context, db *sql.DB, startId int, rowsRet int) ([]RecipeInputOutputInfo, error) {
	query := "SELECT * FROM Recipes_inputs WHERE id >= " + strconv.Itoa(startId)
	if rowsRet > 0 {
		query += " LIMIT " + strconv.Itoa(rowsRet)
	}
	query += ";"
	result, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from db: %w", err)
	}
	var resultRows []RecipeInputOutputInfo
	for result.Next() {
		var row RecipeInputOutputInfo
		err = result.Scan(&row.Id, &row.UsersId, &row.RecipesId, &row.ResourcesId, &row.Amount)
		if err != nil {
			return nil, fmt.Errorf("could not parse data retrieved from db: %w", err)
		}
		resultRows = append(resultRows, row)
	}
	err = result.Err()
	if err != nil {
		return nil, fmt.Errorf("encountered an unexpected error: %w", err)
	}
	return resultRows, nil
}

func InsertRecipesInputs(ctx context.Context, db *sql.DB, data []RecipeInputOutputInfo) (sql.Result, error) {
	query := "INSERT INTO recipes_inputs(users_id, recipes_id, resources_id, amount) VALUES"
	for i, entry := range data {
		if i != 0 {
			query += ","
		}
		query += ` ("` + fmt.Sprint(entry.UsersId) +
			`", ` + fmt.Sprint(entry.RecipesId) +
			`, ` + fmt.Sprint(entry.ResourcesId) +
			`, "` + fmt.Sprint(entry.Amount) + `")`
	}
	query += ";"
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been inserted: %w", err)
	}
	return result, nil
}

func DeleteRecipesInputs(ctx context.Context, db *sql.DB, ids []int) (sql.Result, error) {
	query := "DELETE FROM recipes_inputs WHERE id in ("
	for i, id := range ids {
		if i != 0 {
			query += ","
		}
		query += " " + fmt.Sprint(id)
	}
	query += ");"
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been deleted: %w", err)
	}
	return result, nil
}

func UpdateRecipesInputs(ctx context.Context, db *sql.DB, data []RecipeInputOutputInfo) ([]sql.Result, error) {
	results := []sql.Result{}
	for _, entry := range data {
		query := fmt.Sprintf("UPDATE recipes_inputs SET recipes_id='%d', resources_id=%d, amount='%d' WHERE id=%d;",
			entry.RecipesId, entry.ResourcesId, entry.Amount, entry.Id)
		result, err := db.ExecContext(ctx, query)
		results = append(results, result)
		if err != nil {
			return results, fmt.Errorf("data has not been fully updated: %w", err)
		}
	}
	return results, nil
}

func SelectRecipesOutputs(ctx context.Context, db *sql.DB, startId int, rowsRet int) ([]RecipeInputOutputInfo, error) {
	query := "SELECT * FROM Recipes_outputs WHERE id >= " + strconv.Itoa(startId)
	if rowsRet > 0 {
		query += " LIMIT " + strconv.Itoa(rowsRet)
	}
	query += ";"
	result, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from db: %w", err)
	}
	var resultRows []RecipeInputOutputInfo
	for result.Next() {
		var row RecipeInputOutputInfo
		err = result.Scan(&row.Id, &row.UsersId, &row.RecipesId, &row.ResourcesId, &row.Amount)
		if err != nil {
			return nil, fmt.Errorf("could not parse data retrieved from db: %w", err)
		}
		resultRows = append(resultRows, row)
	}
	err = result.Err()
	if err != nil {
		return nil, fmt.Errorf("encountered an unexpected error: %w", err)
	}
	return resultRows, nil
}

func InsertRecipesOutputs(ctx context.Context, db *sql.DB, data []RecipeInputOutputInfo) (sql.Result, error) {
	query := "INSERT INTO recipes_outputs(users_id, recipes_id, resources_id, amount) VALUES"
	for i, entry := range data {
		if i != 0 {
			query += ","
		}
		query += ` ("` + fmt.Sprint(entry.UsersId) +
			`", ` + fmt.Sprint(entry.RecipesId) +
			`, ` + fmt.Sprint(entry.ResourcesId) +
			`, "` + fmt.Sprint(entry.Amount) + `")`
	}
	query += ";"
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been inserted: %w", err)
	}
	return result, nil
}

func DeleteRecipesOutputs(ctx context.Context, db *sql.DB, ids []int) (sql.Result, error) {
	query := "DELETE FROM recipes_outputs WHERE id in ("
	for i, id := range ids {
		if i != 0 {
			query += ","
		}
		query += " " + fmt.Sprint(id)
	}
	query += ");"
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been deleted: %w", err)
	}
	return result, nil
}

func UpdateRecipesOutputs(ctx context.Context, db *sql.DB, data []RecipeInputOutputInfo) ([]sql.Result, error) {
	results := []sql.Result{}
	for _, entry := range data {
		query := fmt.Sprintf("UPDATE recipes_outputs SET recipes_id='%d', resources_id=%d, amount='%d' WHERE id=%d;",
			entry.RecipesId, entry.ResourcesId, entry.Amount, entry.Id)
		result, err := db.ExecContext(ctx, query)
		results = append(results, result)
		if err != nil {
			return results, fmt.Errorf("data has not been fully updated: %w", err)
		}
	}
	return results, nil
}
