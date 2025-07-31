package application

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DbAddress      string
	ServerPort     uint16
	ServerCertPath string
}

func LoadConfig() Config {
	cfg := Config{
		DbAddress:      "127.0.0.1:3306",
		ServerPort:     3000,
		ServerCertPath: "crud_microservice_secret.pem",
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
	if serverCertPath, exists := os.LookupEnv("CERT"); exists {
		cfg.ServerCertPath = serverCertPath
		fmt.Println("Found server certificate file path", serverCertPath)
	}
	return cfg
}
