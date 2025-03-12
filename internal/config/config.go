package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	const configFileName = ".gatorconfig.json"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error setting home directory: %v", err)
	}

	path := fmt.Sprintf("%s%s%s", homeDir, "/", configFileName)

	return path, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting config filepath: %v", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	js, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	_, err = file.Write(js)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func Read() (Config, error) {
	config := Config{}

	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("error getting filepath: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("error decoding JSON: %v", err)
	}

	return config, nil
}

func (c Config) SetUser(currentUser string) error {
	c.CurrentUserName = currentUser
	err := write(c)
	if err != nil {
		return fmt.Errorf("error setting current username: %v", err)
	}

	return nil
}
