package dbconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"service/backup/databases/client/model"
)

func ListDatabaseConfig(pathFile model.PathFile) {
	fmt.Println("List Databases")

	listDatabases := []model.Database{}

	dataJson, err := os.ReadFile(pathFile.PathDBJson)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	err = json.Unmarshal(dataJson, &listDatabases)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	for _, value := range listDatabases {
		fmt.Println("- " + value.DatabaseName)
	}
}
