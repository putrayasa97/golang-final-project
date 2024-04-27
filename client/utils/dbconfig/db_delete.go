package dbconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"service/backup/databases/client/model"
)

func DeleteDatabaseConfig(dbName string, pathFile model.PathFile) error {
	data, err := os.ReadFile(pathFile.PathDBJson)
	if err != nil {
		return err
	}

	var configs []model.Database
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return err
	}

	var updatedConfigs []model.Database
	for _, config := range configs {
		if config.DatabaseName != dbName {
			updatedConfigs = append(updatedConfigs, config)
		}
	}

	newData, err := json.MarshalIndent(updatedConfigs, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(pathFile.PathDBJson, newData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Data deleted successfully.")
	return nil
}
