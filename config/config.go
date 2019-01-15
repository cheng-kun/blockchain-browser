package config

import (
	"encoding/json"
	"navigator-api/logs"
	"os"
)

const ConfigPath = "./config/config.json"

var config Config

type Config struct {
	NBAILedgerHost string
	DbUsername     string
	DbPwd          string
	DbUrlAndPort   string
	DbSchemaName   string
	DbArgs         string
}

func Init() {
	file, _ := os.Open(ConfigPath)
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Config{}
	err := decoder.Decode(&config)
	if err != nil {
		logs.GetLogger().Fatal("error:", err)
	}
}

func GetConfig() Config {
	return config
}
