package application

import (
	"os"
	"strconv"
)

type Config struct {
	DbAddress  string
	ServerPort uint16
}

func LoadConfig() Config {
	cfg := Config{
		DbAddress:  "127.0.0.1:3306",
		ServerPort: 3000,
	}
	if dbAddr, exists := os.LookupEnv("MYSQL_ADDR"); exists {
		cfg.DbAddress = dbAddr
	}
	if serverPort, exists := os.LookupEnv("PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}
	return cfg
}
