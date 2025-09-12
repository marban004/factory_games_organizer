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

package recipeoutput

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marban004/factory_games_organizer/microservice_logic_crud/model"
)

type MySQLRepo struct {
	DB *sql.DB
}

func (r *MySQLRepo) SelectRecipesOutputsById(ctx context.Context, ids []int, userId int) ([]model.RecipeInputOutputInfo, error) {
	query := "SELECT * FROM recipes_outputs WHERE id in ("
	for i, id := range ids {
		if i != 0 {
			query += ","
		}
		query += " " + fmt.Sprint(id)
	}
	query += ") AND users_id = " + fmt.Sprint(userId) + ";"
	result, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from db: %w", err)
	}
	var resultRows []model.RecipeInputOutputInfo
	for result.Next() {
		var row model.RecipeInputOutputInfo
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

func (r *MySQLRepo) SelectRecipesOutputs(ctx context.Context, startId int, rowsRet int, userId int) ([]model.RecipeInputOutputInfo, error) {
	query := "SELECT * FROM Recipes_outputs WHERE id >= " + fmt.Sprint(startId) + " AND users_id = " + fmt.Sprint(userId)
	if rowsRet > 0 {
		query += " LIMIT " + fmt.Sprint(rowsRet)
	}
	query += ";"
	result, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from db: %w", err)
	}
	var resultRows []model.RecipeInputOutputInfo
	for result.Next() {
		var row model.RecipeInputOutputInfo
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

func (r *MySQLRepo) InsertRecipesOutputs(ctx context.Context, data []model.RecipeInputOutputInfo) (sql.Result, error) {
	query := "INSERT INTO recipes_outputs(users_id, recipes_id, resources_id, amount) VALUES"
	for i, entry := range data {
		if i != 0 {
			query += ","
		}
		query += ` (` + fmt.Sprint(entry.UsersId) +
			`, ` + fmt.Sprint(entry.RecipesId) +
			`, ` + fmt.Sprint(entry.ResourcesId) +
			`, ` + fmt.Sprint(entry.Amount) + `)`
	}
	query += ";"
	result, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been inserted: %w", err)
	}
	return result, nil
}

func (r *MySQLRepo) DeleteRecipesOutputs(ctx context.Context, ids []int, userId int) (sql.Result, error) {
	query := "DELETE FROM recipes_outputs WHERE id in ("
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

func (r *MySQLRepo) DeleteRecipesOutputsByUserId(ctx context.Context, transaction *sql.Tx, userId int) (sql.Result, error) {
	query := "DELETE FROM recipes_outputs WHERE users_id = " + fmt.Sprint(userId) + ";"
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

func (r *MySQLRepo) UpdateRecipesOutputs(ctx context.Context, data []model.RecipeInputOutputInfo) ([]sql.Result, error) {
	results := []sql.Result{}
	transaction, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return results, fmt.Errorf("data has not been updated: %w", err)
	}
	for _, entry := range data {
		query := fmt.Sprintf("UPDATE recipes_outputs SET recipes_id='%d', resources_id=%d, amount='%d' WHERE id=%d and users_id=%d;",
			entry.RecipesId, entry.ResourcesId, entry.Amount, entry.Id, entry.UsersId)
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
