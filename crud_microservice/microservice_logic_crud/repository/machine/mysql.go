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

package machine

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marban004/factory_games_organizer/microservice_logic_crud/model"
)

type MySQLRepo struct {
	DB *sql.DB
}

func (r *MySQLRepo) SelectMachinesById(ctx context.Context, ids []int, userId int) ([]model.MachineInfo, error) {
	query := "SELECT * FROM machines WHERE id in ("
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
	var resultRows []model.MachineInfo
	for result.Next() {
		var row model.MachineInfo
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

func (r *MySQLRepo) SelectMachines(ctx context.Context, startId int, rowsRet int, userId int) ([]model.MachineInfo, error) {
	query := "SELECT * FROM machines WHERE id >= " + fmt.Sprint(startId) + " AND users_id = " + fmt.Sprint(userId)
	if rowsRet > 0 {
		query += " LIMIT " + fmt.Sprint(rowsRet)
	}
	query += ";"
	result, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from db: %w", err)
	}
	var resultRows []model.MachineInfo
	for result.Next() {
		var row model.MachineInfo
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

func (r *MySQLRepo) InsertMachines(ctx context.Context, data []model.MachineInfo, userId uint) (sql.Result, error) {
	query := "INSERT INTO machines(name, users_id, inputs_solid, inputs_liquid, outputs_solid, outputs_liquid, speed, power_consumption_kw, default_choice) VALUES"
	for i, entry := range data {
		if i != 0 {
			query += ","
		}
		query += ` ("` + entry.Name +
			`", ` + fmt.Sprint(userId) +
			`, ` + fmt.Sprint(entry.InputsSolid) +
			`, ` + fmt.Sprint(entry.InputsLiquid) +
			`, ` + fmt.Sprint(entry.OutputsSolid) +
			`, ` + fmt.Sprint(entry.OutputsLiquid) +
			`, ` + fmt.Sprint(entry.Speed) +
			`, ` + fmt.Sprint(entry.PowerConsumptionKw) +
			`, ` + fmt.Sprint(entry.DefaultChoice) + `)`
	}
	query += ";"
	result, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been inserted: %w", err)
	}
	return result, nil
}

func (r *MySQLRepo) DeleteMachines(ctx context.Context, ids []int, userId int) (sql.Result, error) {
	query := "DELETE FROM machines WHERE id in ("
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

func (r *MySQLRepo) DeleteMachinesByUserId(ctx context.Context, transaction *sql.Tx, userId int) (sql.Result, error) {
	query := "DELETE FROM machines WHERE users_id = " + fmt.Sprint(userId) + ";"
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

func (r *MySQLRepo) UpdateMachines(ctx context.Context, data []model.MachineInfo, userId uint) ([]sql.Result, error) {
	results := []sql.Result{}
	transaction, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return results, fmt.Errorf("data has not been updated: %w", err)
	}
	for _, entry := range data {
		query := fmt.Sprintf("UPDATE machines SET name='%s', inputs_solid=%d, inputs_liquid=%d, outputs_solid=%d, outputs_liquid=%d, speed=%f, power_consumption_kw=%d, default_choice=%d WHERE id=%d and users_id=%d;",
			entry.Name, entry.InputsSolid, entry.InputsLiquid, entry.OutputsSolid, entry.OutputsLiquid, entry.Speed, entry.PowerConsumptionKw, entry.DefaultChoice, entry.Id, userId)
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
