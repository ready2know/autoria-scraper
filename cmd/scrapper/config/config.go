package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type AppConfig struct {
	Postgres struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
	} `json:"postgres"`
}

func LoadConfig(appConfig *AppConfig) {
	file, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&appConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
}
