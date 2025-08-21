package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/marban004/factory_games_organizer/microservice_logic_dispatcher/model"
)

type MySQLRepo struct {
	DB *sql.DB
}

func (r *MySQLRepo) SelectUserByLogin(ctx context.Context, login string) (model.UserInfo, error) {
	user := model.UserInfo{}
	query := fmt.Sprintf(`SELECT * FROM users where login = "%s"`, strings.ToLower(login))
	err := r.DB.QueryRowContext(ctx, query).Scan(&user.UserId, &user.UserLogin, &user.UserPasswdHash)
	if err != nil {
		return user, fmt.Errorf("could not retrive information from database: %w", err)
	}
	return user, nil
}

// func (r *MySQLRepo) SelectUserById(ctx context.Context, userId uint) (model.UserInfo, error) {
// 	user := model.UserInfo{}
// 	query := fmt.Sprintf(`SELECT * FROM users where id = %d`, userId)
// 	err := r.DB.QueryRowContext(ctx, query).Scan(&user.UserId, &user.UserLogin, &user.UserPasswdHash)
// 	if err != nil {
// 		return user, fmt.Errorf("could not retrive information from database: %w", err)
// 	}
// 	return user, nil
// }

// put in User repository, only creates does not verify validity of login and password
func (r *MySQLRepo) CreateUser(ctx context.Context, user model.UserInfo) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT INTO users VALUES (null, "%s", "%s")`, strings.ToLower(user.UserLogin), user.UserPasswdHash)
	result, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not insert new user: %w", err)
	}
	return result, nil
}

// put in User repository, does not verify validity of login and password, replace UserId in user with id from jwt
func (r *MySQLRepo) UpdateUser(ctx context.Context, user model.UserInfo) (sql.Result, error) {
	query := `UPDATE users SET `
	if len(user.UserLogin) > 0 {
		query += fmt.Sprintf(`login = "%s" `, strings.ToLower(user.UserLogin))
	}
	if len(user.UserPasswdHash) > 0 {
		query += fmt.Sprintf(`passwdhash = "%s" `, user.UserPasswdHash)
	}
	query += fmt.Sprintf(`WHERE id = %d`, user.UserId)
	result, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("user has not been updated: %w", err)
	}
	return result, nil
}

func (r *MySQLRepo) DeleteUser(ctx context.Context, userId uint) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM users WHERE id = %d", userId)
	result, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("data has not been deleted: %w", err)
	}
	return result, nil
}
