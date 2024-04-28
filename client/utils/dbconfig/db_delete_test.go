package dbconfig_test

import (
	"encoding/json"
	"os"
	"service/backup/databases/client/model"
	"service/backup/databases/client/utils/dbconfig"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteDatabaseConfig(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testdb.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("[]")
	if err != nil {
		t.Fatalf("Gagal menulis ke file sementara: %v", err)
	}
	var configs []model.Database
	config1 := []string{"dbname", "localhost", "5432", "user", "password"}
	config2 := []string{"dbname1", "localhost", "2222", "user2", "password2"}
	pathFile := model.PathFile{PathDBJson: tempFile.Name()}

	assert.Equal(t, len(configs), 0)

	err = dbconfig.AddDatabaseConfig(config1, pathFile)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	err = dbconfig.AddDatabaseConfig(config2, pathFile)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	data1, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temporary file: %v", err)
	}

	err = json.Unmarshal(data1, &configs)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	assert.Equal(t, len(configs), 2)

	err = dbconfig.DeleteDatabaseConfig("dbname1", pathFile)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	data2, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temporary file: %v", err)
	}

	err = json.Unmarshal(data2, &configs)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	assert.Equal(t, len(configs), 1)
}
