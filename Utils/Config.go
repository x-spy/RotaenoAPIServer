package Utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIVersion          string              `json:"apiVersion"`
	ServerPort          string              `json:"serverPort"`
	FilePath            string              `json:"filePath"`
	PurchaseInformation PurchaseInformation `json:"purchaseInformation"`
	AllowedObjectID     []string            `json:"allowedObjectID"`
}

type PurchaseInformation struct {
	BaseObjectID string            `json:"baseObjectID"`
	Keys         map[string]string `json:"keys"`
}

func InitConfig() error {
	executablePath, _ := os.Executable()
	rootPath := filepath.Dir(executablePath)
	_, err := os.Stat(filepath.Join(rootPath, "config.json"))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Couldn't find config file. Creating new config...")
			err := createDefaultConfig()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func createDefaultConfig() error {
	var defaultConfig Config
	defaultConfig.APIVersion = "1.1"
	defaultConfig.ServerPort = "8080"
	defaultConfig.FilePath = "files"
	defaultConfig.PurchaseInformation.BaseObjectID = "base-object-id"
	defaultConfig.AllowedObjectID = make([]string, 0)

	executablePath, _ := os.Executable()
	rootPath := filepath.Dir(executablePath)
	configFile, err := os.Create(filepath.Join(rootPath, "config.json"))
	if err != nil {
		return err
	}
	defaultConfigBytes, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		return err
	}
	_, err = configFile.Write(defaultConfigBytes)
	if err != nil {
		return err
	}

	return nil
}

func GetConfig() (Config, error) {
	executablePath, _ := os.Executable()
	rootPath := filepath.Dir(executablePath)
	var config Config

	// Read Config
	fmt.Println("Reading config")
	configPath := filepath.Join(rootPath, "config.json")
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil

}
