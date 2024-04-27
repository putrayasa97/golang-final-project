package dbconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"service/backup/databases/client/model"
)

func AddDatabaseConfig(config []string, pathFile model.PathFile) error {
	data, err := os.ReadFile(pathFile.PathDBJson)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	var configs []model.Database
	err = json.Unmarshal(data, &configs)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}

	newConfig := model.Database{
		DatabaseName: config[0],
		DBHost:       config[1],
		DBPort:       config[2],
		DBUsername:   config[3],
		DBPassword:   config[4],
	}
	configs = append(configs, newConfig)

	newData, err := json.MarshalIndent(configs, "", "    ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err
	}

	err = os.WriteFile(pathFile.PathDBJson, newData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	fmt.Println("Data added successfully.")

	return nil
}
