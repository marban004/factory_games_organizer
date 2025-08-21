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
	UsersMicroservicesAddresses      []string
	CrudMicroservicesAddresses       []string
	CalculatorMicroservicesAddresses []string
}

func LoadConfig() Config {
	cfg := Config{
		DbAddress:                        "127.0.0.1:3306",
		ServerPort:                       3000,
		ServerSecretPath:                 "dispatcher_microservice_secret.pem",
		ServerCertPath:                   "dispatcher_microservice_cert.crt",
		UsersMicroservicesAddresses:      []string{"127.0.0.1:8082"},
		CrudMicroservicesAddresses:       []string{"127.0.0.1:8081"},
		CalculatorMicroservicesAddresses: []string{"127.0.0.1:8080"},
	}
	if dbAddr, exists := os.LookupEnv("MYSQL_ADDR"); exists {
		cfg.DbAddress = dbAddr
		fmt.Println("Found db address:", cfg.DbAddress)
	}
	if serverPort, exists := os.LookupEnv("PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
			fmt.Println("Found server port:", cfg.ServerPort)
		}
	}
	if serverSecretPath, exists := os.LookupEnv("SECRET"); exists {
		cfg.ServerSecretPath = serverSecretPath
		fmt.Println("Found secret key file path", serverSecretPath)
	}
	if serverCertPath, exists := os.LookupEnv("CERT"); exists {
		cfg.ServerCertPath = serverCertPath
		fmt.Println("Found certificate file path", serverCertPath)
	}
	if usersMicroservicesAddresses, exists := os.LookupEnv("USERS"); exists {
		cfg.UsersMicroservicesAddresses = strings.Split(usersMicroservicesAddresses, ",")
		fmt.Println("Found Users microservice URL list", cfg.UsersMicroservicesAddresses)
	}
	if crudMicroservicesAddresses, exists := os.LookupEnv("CRUD"); exists {
		cfg.CrudMicroservicesAddresses = strings.Split(crudMicroservicesAddresses, ",")
		fmt.Println("Found Users microservice URL list", cfg.CrudMicroservicesAddresses)
	}
	if calculatorMicroservicesAddresses, exists := os.LookupEnv("CALCULATOR"); exists {
		cfg.CalculatorMicroservicesAddresses = strings.Split(calculatorMicroservicesAddresses, ",")
		fmt.Println("Found Users microservice URL list", cfg.CalculatorMicroservicesAddresses)
	}
	return cfg
}
