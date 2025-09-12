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
	"os"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/marban004/factory_games_organizer/prototypes"
	"github.com/stretchr/testify/suite"
)

type UsersPrototypeIntegrationTestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestUsersPrototypeIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &UsersPrototypeIntegrationTestSuite{})
}

func (upits *UsersPrototypeIntegrationTestSuite) SetupSuite() {
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = "pH082C./"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "users_data"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		upits.FailNowf("unable to connect to database, error: %s", err.Error())
	}
	upits.db = db
	setupDatabaseSchema(upits)
	cleanDatabaseTables(upits)
}

func (upits *UsersPrototypeIntegrationTestSuite) TearDownSuite() {
	teardownDB(upits)
}

func (upits *UsersPrototypeIntegrationTestSuite) TearDownTest() {
	cleanDatabaseTables(upits)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserPositive() {
	login := "mat"
	password := "test_password"
	userId, err := prototypes.VerifyUser(context.Background(), upits.db, login, password)
	upits.Nil(err)
	upits.EqualValues(1, userId, "actual value differs from expected")
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserNegative() {
	login := "mat"
	password := "test_password_wrong"
	userId, err := prototypes.VerifyUser(context.Background(), upits.db, login, password)
	upits.Nil(err)
	upits.EqualValues(0, userId, "actual value differs from expected")
}

func (upits *UsersPrototypeIntegrationTestSuite) TestCreateUserPositive() {
	user := prototypes.User{
		UserId:         2,
		UserLogin:      "mat2",
		UserPasswdHash: "",
	}
	password := "gyY*89$$"
	user.UserPasswdHash, _ = prototypes.GeneratePasswordHash(password)
	result, err := prototypes.CreateUser(context.Background(), upits.db, user)
	upits.Nil(err)
	rowsChanged, err := result.RowsAffected()
	upits.Nil(err)
	upits.Equal(int64(1), rowsChanged, "The number of changed rows differs from expected")
	userId, err := prototypes.VerifyUser(context.Background(), upits.db, user.UserLogin, password)
	upits.Nil(err)
	upits.NotEqualValues(0, userId, "actual value differs from expected")
}

func (upits *UsersPrototypeIntegrationTestSuite) TestCreateUserNegative() {
	user := prototypes.User{
		UserId:         2,
		UserLogin:      "mat",
		UserPasswdHash: "",
	}
	password := "gyY*89$$"
	user.UserPasswdHash, _ = prototypes.GeneratePasswordHash(password)
	result, err := prototypes.CreateUser(context.Background(), upits.db, user)
	upits.NotNil(err)
	upits.Nil(result)
	userId, err := prototypes.VerifyUser(context.Background(), upits.db, user.UserLogin, password)
	upits.Nil(err)
	upits.EqualValues(0, userId, "actual value differs from expected")
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserLoginPositive() {
	login := "mat"
	valid, err := prototypes.VerifyUserLogin(login)
	upits.Nil(err)
	upits.Equal(true, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserLoginNegative1() {
	login := "mt"
	valid, err := prototypes.VerifyUserLogin(login)
	upits.Nil(err)
	upits.Equal(false, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserLoginNegative2() {
	login := "mt`;"
	valid, err := prototypes.VerifyUserLogin(login)
	upits.Nil(err)
	upits.Equal(false, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserLoginNegative3() {
	login := `mt";SELECT * FROM users;`
	valid, err := prototypes.VerifyUserLogin(login)
	upits.Nil(err)
	upits.Equal(false, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserPasswordPositive() {
	passwd := "gyY*89$$"
	valid, err := prototypes.VerifyUserPassword(passwd)
	upits.Nil(err)
	upits.Equal(true, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserPasswordNegative1() {
	passwd := "gyY*89$"
	valid, err := prototypes.VerifyUserPassword(passwd)
	upits.Nil(err)
	upits.Equal(false, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserPasswordNegative2() {
	passwd := "gyuu*89$$"
	valid, err := prototypes.VerifyUserPassword(passwd)
	upits.Nil(err)
	upits.Equal(false, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserPasswordNegative3() {
	passwd := "GTTF*89$"
	valid, err := prototypes.VerifyUserPassword(passwd)
	upits.Nil(err)
	upits.Equal(false, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserPasswordNegative4() {
	passwd := "GTTFi89iii"
	valid, err := prototypes.VerifyUserPassword(passwd)
	upits.Nil(err)
	upits.Equal(false, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserPasswordNegative5() {
	passwd := "gyY*jjjjj$"
	valid, err := prototypes.VerifyUserPassword(passwd)
	upits.Nil(err)
	upits.Equal(false, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserPasswordNegative6() {
	passwd := "gyY*89$`;"
	valid, err := prototypes.VerifyUserPassword(passwd)
	upits.Nil(err)
	upits.Equal(false, valid)
}

func (upits *UsersPrototypeIntegrationTestSuite) TestVerifyUserPasswordNegative7() {
	passwd := "gyY*89$SElect"
	valid, err := prototypes.VerifyUserPassword(passwd)
	upits.Nil(err)
	upits.Equal(false, valid)
}

func setupDatabaseSchema(upits *UsersPrototypeIntegrationTestSuite) {
	upits.T().Log("deleting previous schema")
	_, err := upits.db.Exec(`DROP DATABASE IF EXISTS users_test`)
	if err != nil {
		upits.FailNowf("unable to drop previous database", err.Error())
	}
	upits.T().Log("setting up database schema")
	_, err = upits.db.Exec(`CREATE DATABASE users_test`)
	if err != nil {
		upits.FailNowf("unable to create database", err.Error())
	}

	_, err = upits.db.Exec(`USE users_test`)
	if err != nil {
		upits.FailNowf("unable to switch currently used database", err.Error())
	}

	contents, err := os.ReadFile("users_schema_mysql.sql")
	if err != nil {
		upits.FailNowf("unable to read schema sql file", err.Error())
	}
	commands := strings.Split(string(contents), ";")
	for _, command := range commands {
		if len(command) <= 0 {
			break
		}
		_, err = upits.db.Exec(command)
		if err != nil {
			upits.FailNowf("unable to setup database schema", err.Error())
		}
	}
}

func cleanDatabaseTables(upits *UsersPrototypeIntegrationTestSuite) {
	upits.T().Log("cleaning database tables")
	contents, err := os.ReadFile("users_data_mysql.sql")
	if err != nil {
		upits.FailNowf("unable to read data sql file", err.Error())
	}

	commands := strings.Split(string(contents), ";")
	for _, command := range commands {
		if len(command) <= 0 {
			break
		}
		_, err = upits.db.Exec(command)
		if err != nil {
			upits.FailNowf("unable to setup database contents", err.Error())
		}
	}
}

func teardownDB(upits *UsersPrototypeIntegrationTestSuite) {
	_, err := upits.db.Exec("DROP DATABASE users_test")
	if err != nil {
		upits.FailNowf("failed to drop test db", err.Error())
	}
	err = upits.db.Close()
	if err != nil {
		upits.FailNowf("failed to close test db", err.Error())
	}
}
