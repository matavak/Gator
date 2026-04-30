package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Could not get user home Dir:%v", err)
	}
	configPath := filepath.Join(userHomeDir + "/.gatorconfig.json")

	return configPath, nil
}

func Read() Config {
	configPath, err := getConfigFilePath()
	if err != nil {
		log.Fatal("could not get user home dir")
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Now let's unmarshall the data into `payload`
	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}

func (cfg *Config) SetUser(userName string) {
	cfg.CurrentUserName = userName
}

func (cfg *Config) write() {
	configPath, err := getConfigFilePath()
	if err != nil {
		log.Fatal("could not get user home dir")
	}
	byteValue, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.WriteFile(configPath, byteValue, 0o644)
}
