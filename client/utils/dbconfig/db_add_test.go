package dbconfig_test

import (
	"encoding/json"
	"os"
	"service/backup/databases/client/model"
	"service/backup/databases/client/utils/dbconfig"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDatabaseConfig(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testdb.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("[]")
	if err != nil {
		t.Fatalf("Gagal menulis ke file sementara: %v", err)
	}

	config := []string{"dbname", "localhost", "5432", "user", "password"}
	pathFile := model.PathFile{PathDBJson: tempFile.Name()}

	err = dbconfig.AddDatabaseConfig(config, pathFile)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temporary file: %v", err)
	}

	var configs []model.Database
	err = json.Unmarshal(data, &configs)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	assert.Equal(t, configs[0].DatabaseName, config[0])
	assert.Equal(t, configs[0].DBHost, config[1])
	assert.Equal(t, configs[0].DBPort, config[2])
	assert.Equal(t, configs[0].DBUsername, config[3])
	assert.Equal(t, configs[0].DBPassword, config[4])
}
