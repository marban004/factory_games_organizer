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

package application

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DbAddress                        string
	ServerPort                       uint16
	ServerSecretPath                 string
	ServerCertPath                   string
	Host                             string
	ApiStatsFile                     string
	DumpStats                        bool
	TrackerCapacity                  uint64
	TrackerTimePeriod                int64
	UsersMicroservicesAddresses      []string
	CrudMicroservicesAddresses       []string
	CalculatorMicroservicesAddresses []string
}

func LoadConfig() Config {
	cfg := Config{
		// DbAddress:                        "127.0.0.1:3306",
		ServerPort:                       3000,
		TrackerCapacity:                  1440,
		TrackerTimePeriod:                60000,
		ServerSecretPath:                 "dispatcher_microservice_secret.pem",
		ServerCertPath:                   "dispatcher_microservice_cert.crt",
		Host:                             "localhost",
		ApiStatsFile:                     "",
		DumpStats:                        true,
		UsersMicroservicesAddresses:      []string{"127.0.0.1:8082"},
		CrudMicroservicesAddresses:       []string{"127.0.0.1:8081"},
		CalculatorMicroservicesAddresses: []string{"127.0.0.1:8080"},
	}
	// if dbAddr, exists := os.LookupEnv("MYSQL_ADDR"); exists {
	// 	cfg.DbAddress = dbAddr
	// 	fmt.Println("Found db address:", cfg.DbAddress)
	// }
	if serverPort, exists := os.LookupEnv("PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
			fmt.Println("Found server port:", cfg.ServerPort)
		}
	}
	if serverSecretPath, exists := os.LookupEnv("SECRET"); exists {
		cfg.ServerSecretPath = serverSecretPath
		fmt.Println("Found secret key file path:", serverSecretPath)
	}
	if serverCertPath, exists := os.LookupEnv("CERT"); exists {
		cfg.ServerCertPath = serverCertPath
		fmt.Println("Found certificate file path:", serverCertPath)
	}
	if usersMicroservicesAddresses, exists := os.LookupEnv("USERS"); exists {
		cfg.UsersMicroservicesAddresses = strings.Split(usersMicroservicesAddresses, ",")
		fmt.Println("Found Users microservice URL list:", cfg.UsersMicroservicesAddresses)
	}
	if crudMicroservicesAddresses, exists := os.LookupEnv("CRUD"); exists {
		cfg.CrudMicroservicesAddresses = strings.Split(crudMicroservicesAddresses, ",")
		fmt.Println("Found Crud microservice URL list:", cfg.CrudMicroservicesAddresses)
	}
	if calculatorMicroservicesAddresses, exists := os.LookupEnv("CALCULATOR"); exists {
		cfg.CalculatorMicroservicesAddresses = strings.Split(calculatorMicroservicesAddresses, ",")
		fmt.Println("Found Calculator microservice URL list:", cfg.CalculatorMicroservicesAddresses)
	}
	if host, exists := os.LookupEnv("HOST"); exists {
		cfg.Host = host
		fmt.Println("Found host server address:", host)
	}
	if trackerCapacity, exists := os.LookupEnv("TRCAP"); exists {
		if trackerCapacityUint, err := strconv.ParseUint(trackerCapacity, 10, 64); err == nil {
			cfg.TrackerCapacity = trackerCapacityUint
			fmt.Println("Found stat tracker capacity:", cfg.TrackerCapacity)
		}
	}
	if trackerPeriod, exists := os.LookupEnv("TRPER"); exists {
		if trackerPeriodInt, err := strconv.ParseInt(trackerPeriod, 10, 64); err == nil {
			cfg.TrackerTimePeriod = trackerPeriodInt
			fmt.Println("Found stat tracker measure period:", cfg.TrackerTimePeriod)
		}
	}
	if outFile, exitst := os.LookupEnv("OUT"); exitst {
		cfg.ApiStatsFile = outFile
		fmt.Println("Found name for api stats file:", outFile)
	}
	if dumpStats, exists := os.LookupEnv("DUMPSTATS"); exists {
		if temp, err := strconv.ParseInt(dumpStats, 10, 8); err == nil {
			if temp > 0 {
				cfg.DumpStats = true
				fmt.Println("Api stats set to be dumped after stopping server")
			}
		}
	}
	return cfg
}
