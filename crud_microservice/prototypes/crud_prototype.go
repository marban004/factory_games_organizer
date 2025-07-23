package prototypes

import (
	"context"
	"database/sql"
	"fmt"
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

type JSONInput struct {
	MachinesList  []MachineInfo
	ResourcesList []ResourceInfo
}

//todo: add selects for all tables (select everything with a upper limit on rows selected, if no limit provided select everything), if enough time allow selects of rows in id list

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
