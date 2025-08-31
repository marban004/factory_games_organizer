package application

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DbAddress         string
	ServerPort        uint16
	ServerSecretPath  string
	ServerCertPath    string
	Host              string
	ApiStatsFile      string
	DumpStats         bool
	TrackerCapacity   uint64
	TrackerTimePeriod int64
}

func LoadConfig() Config {
	cfg := Config{
		DbAddress:         "127.0.0.1:3306",
		ServerPort:        3000,
		ServerSecretPath:  "users_microservice_secret.pem",
		ServerCertPath:    "users_microservice_cert.crt",
		Host:              "localhost",
		TrackerCapacity:   1440,
		TrackerTimePeriod: 60000,
		ApiStatsFile:      "",
		DumpStats:         true,
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
