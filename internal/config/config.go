package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	directory, err := configFilePathGet()
	if err != nil {
		var error Config
		return error, err
	}

	file, err := os.Open(directory)
	if err != nil {
		var error Config
		return error, fmt.Errorf("error: %v", err)
	}
	defer file.Close()

	filebytes, err := io.ReadAll(file)
	if err != nil {
		var error Config
		return error, fmt.Errorf("error: %v", err)
	}

	var jsonConfig Config
	if err := json.Unmarshal(filebytes, &jsonConfig); err != nil {
		var error Config
		return error, fmt.Errorf("error: %v", err)
	}

	return jsonConfig, nil

}

func (cfg *Config) SetUser(user string) error {
	currentConfig := cfg
	currentConfig.Current_user_name = user

	jsonData, err := json.Marshal(currentConfig)
	if err != nil {
		return err
	}

	directory, err := configFilePathGet()
	if err != nil {
		return err
	}

	err = os.WriteFile(directory, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func configFilePathGet() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("error: %v", err)
	}
	directory := homeDir + "/" + ".gatorconfig.json"

	return directory, nil
}
