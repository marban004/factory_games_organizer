package machinerecipe

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marban004/factory_games_organizer/microservice_logic_crud/model"
)

type MySQLRepo struct {
	DB *sql.DB
}

func (r *MySQLRepo) SelectMachinesRecipesById(ctx context.Context, ids []int, userId int) ([]model.MachinesRecipesInfo, error) {
	query := "SELECT * FROM machines_recipes WHERE id in ("
	for i, id := range ids {
		if i != 0 {
			query += ","
		}
		query += " " + fmt.Sprint(id)
	}
	query += ") AND users_id = " + fmt.Sprint(userId) + ";"
	result, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from r.DB: %w", err)
	}
	var resultRows []model.MachinesRecipesInfo
	for result.Next() {
		var row model.MachinesRecipesInfo
		err = result.Scan(&row.Id, &row.UsersId, &row.RecipesId, &row.MachinesId)
		if err != nil {
			return nil, fmt.Errorf("could not parse data retrieved from r.DB: %w", err)
		}
		resultRows = append(resultRows, row)
	}
	err = result.Err()
	if err != nil {
		return nil, fmt.Errorf("encountered an unexpected error: %w", err)
	}
	return resultRows, nil
}

func (r *MySQLRepo) SelectMachinesRecipes(ctx context.Context, startId int, rowsRet int, userId int) ([]model.MachinesRecipesInfo, error) {
	query := "SELECT * FROM machines_recipes WHERE id >= " + fmt.Sprint(startId) + " AND users_id = " + fmt.Sprint(userId)
	if rowsRet > 0 {
		query += " LIMIT " + fmt.Sprint(rowsRet)
	}
	query += ";"
	result, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from r.DB: %w", err)
	}
	var resultRows []model.MachinesRecipesInfo
	for result.Next() {
		var row model.MachinesRecipesInfo
		err = result.Scan(&row.Id, &row.UsersId, &row.RecipesId, &row.MachinesId)
		if err != nil {
			return nil, fmt.Errorf("could not parse data retrieved from r.DB: %w", err)
		}
		resultRows = append(resultRows, row)
	}
	err = result.Err()
	if err != nil {
		return nil, fmt.Errorf("encountered an unexpected error: %w", err)
	}
	return resultRows, nil
}

func (r *MySQLRepo) InsertMachinesRecipes(ctx context.Context, data []model.MachinesRecipesInfo) (sql.Result, error) {
	query := "INSERT INTO machines_recipes(users_id, recipes_id, machines_id) VALUES"
	i := 0
	for _, entry := range data {
		err := r.verifyRecipeMachineIntegrity(ctx, entry.RecipesId, entry.MachinesId, entry.UsersId)
		if err.Error() == sql.ErrNoRows.Error() {
			continue
		}
		if i != 0 {
			query += ","
		}
		i++
		query += ` (` + fmt.Sprint(entry.UsersId) +
			`, ` + fmt.Sprint(entry.RecipesId) +
			`, ` + fmt.Sprint(entry.MachinesId) + `)`
	}
	query += ";"
	result, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been inserted: %w", err)
	}
	return result, nil
}

func (r *MySQLRepo) DeleteMachinesRecipes(ctx context.Context, ids []int, userId int) (sql.Result, error) {
	query := "DELETE FROM machines_recipes WHERE id in ("
	for i, id := range ids {
		if i != 0 {
			query += ","
		}
		query += " " + fmt.Sprint(id)
	}
	query += ") and users_id = " + fmt.Sprint(userId) + ";"
	result, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been deleted: %w", err)
	}
	return result, nil
}

func (r *MySQLRepo) DeleteMachinesRecipesByUserId(ctx context.Context, transaction *sql.Tx, userId int) (sql.Result, error) {
	query := "DELETE FROM machines_recipes WHERE users_id = " + fmt.Sprint(userId) + ";"
	result, err := transaction.ExecContext(ctx, query)
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			return nil, fmt.Errorf("could not rollback changes: %w", rollbackErr)
		}
		return nil, fmt.Errorf("an error occurred, transaction has been rolled back: %w", err)
	}
	return result, nil
}

func (r *MySQLRepo) UpdateMachinesRecipes(ctx context.Context, data []model.MachinesRecipesInfo) ([]sql.Result, error) {
	results := []sql.Result{}
	transaction, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return results, fmt.Errorf("data has not been updated: %w", err)
	}
	for _, entry := range data {
		err := r.verifyRecipeMachineIntegrity(ctx, entry.RecipesId, entry.MachinesId, entry.UsersId)
		if err.Error() == sql.ErrNoRows.Error() {
			continue
		}
		query := fmt.Sprintf("UPDATE machines_recipes SET recipes_id='%d', machines_id=%d WHERE id=%d and users_id=%d;",
			entry.RecipesId, entry.MachinesId, entry.Id, entry.UsersId)
		result, err := transaction.ExecContext(ctx, query)
		results = append(results, result)
		if err != nil {
			rollbackErr := transaction.Rollback()
			if rollbackErr != nil {
				return results, fmt.Errorf("could not rollback transaction: %w", rollbackErr)
			}
			return results, fmt.Errorf("data has not been updated: %w", err)
		}
	}
	err = transaction.Commit()
	if err != nil {
		return results, fmt.Errorf("data has not been updated: %w", err)
	}
	return results, nil
}

func (r *MySQLRepo) verifyRecipeMachineIntegrity(ctx context.Context, recipeId uint, machineId uint, userId uint) error {
	query := fmt.Sprintf(`select * from machines where inputs_liquid >= 
	(select count(*) from recipes r 
	left join recipes_inputs ri on r.id = ri.recipes_id 
	left join resources rs on ri.resources_id = rs.id 
	where r.id = %[1]d and users_id = %[3]d and rs.liquid = 1) 
and outputs_liquid >= 
	(select count(*) from recipes r 
	left join recipes_outputs ro on r.id = ro.recipes_id 
	left join resources rs on ro.resources_id = rs.id 
	where r.id = %[1]d and users_id = %[3]d and rs.liquid = 1)
and inputs_solid >= (select case when count(*) > 0 then 1 else 0 end from recipes r 
	left join recipes_inputs ri on r.id = ri.recipes_id 
	left join resources rs on ri.resources_id = rs.id 
	where r.id = %[1]d and users_id = %[3]d and rs.liquid = 0)
and outputs_solid >= 
	(select case when count(*) > 0 then 1 else 0 end from recipes r 
	left join recipes_outputs ro on r.id = ro.recipes_id 
	left join resources rs on ro.resources_id = rs.id 
	where r.id = %[1]d and users_id = %[3]d and rs.liquid = 0)
and id = %[2]d and users_id = %[3]d;`, recipeId, machineId, userId)
	_, err := r.DB.QueryContext(ctx, query)
	if err.Error() == sql.ErrNoRows.Error() {
		return err
	}
	if err != nil {
		return fmt.Errorf("data integrity has not been sucessfully verified: %w", err)
	}
	return nil
}
