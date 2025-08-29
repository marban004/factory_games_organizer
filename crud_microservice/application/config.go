package application

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DbAddress        string
	ServerPort       uint16
	ServerSecretPath string
	ServerCertPath   string
	Host             string
}

func LoadConfig() Config {
	cfg := Config{
		DbAddress:        "127.0.0.1:3306",
		ServerPort:       3000,
		ServerSecretPath: "crud_microservice_secret.pem",
		ServerCertPath:   "crud_microservice_cert.crt",
		Host:             "localhost",
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
		fmt.Println("Found secret key file path:", serverSecretPath)
	}
	if serverCertPath, exists := os.LookupEnv("CERT"); exists {
		cfg.ServerCertPath = serverCertPath
		fmt.Println("Found certificate file path:", serverCertPath)
	}
	if host, exists := os.LookupEnv("HOST"); exists {
		cfg.Host = host
		fmt.Println("Found host server address:", host)
	}
	return cfg
}
