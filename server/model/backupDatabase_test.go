package model_test

import (
	"fmt"
	"service/backup/databases/server/config"
	"service/backup/databases/server/model"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load("../.env.test")
	if err != nil {
		fmt.Println(err)
	}
	config.OpenDB()
}

var databaseBackup model.BckpDatabase

func TestCreateBackup(t *testing.T) {
	InitEnv()

	databaseBackup.DatabaseName = "database_test"
	databaseBackup.FileName = "database_test.sql.zip"
	databaseBackup.FilePath = `C:\Users\comms\Downloads\books.sql`
	databaseBackup.CreatedAt = time.Now()
	databaseBackup.UpdatedAt = time.Now()

	err := databaseBackup.Create(config.Mysql.DB)
	assert.Nil(t, err, "cek error create book")
	assert.NotEmpty(t, databaseBackup.ID, "cek id is not empty")

}
func TestCheckLatestBackup(t *testing.T) {
	InitEnv()

	latestDataBackup, err := databaseBackup.LatestBackup(config.Mysql.DB)
	assert.Nil(t, err, "error create book")
	assert.NotEmpty(t, latestDataBackup, "id is not empty")

	found := false
	for _, data := range latestDataBackup {
		if data.ID == databaseBackup.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "inserted data not found in latest backup")

}

func TestHistoryBackupByName(t *testing.T) {
	InitEnv()

	latestDataBackup, err := databaseBackup.LatestBackup(config.Mysql.DB)
	assert.Nil(t, err, "error create book")
	assert.NotEmpty(t, latestDataBackup, "id is not empty")

	found := false
	for _, data := range latestDataBackup {
		if data.ID == databaseBackup.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "inserted data not found in latest backup")

}

func TestBackupById(t *testing.T) {
	InitEnv()

	dataBackupById, err := databaseBackup.GetBackupById(config.Mysql.DB)
	assert.Nil(t, err, "error create book")
	assert.NotEmpty(t, dataBackupById, "id is not empty")
	assert.Equal(t, dataBackupById.ID, databaseBackup.ID, "Id not same")
}
