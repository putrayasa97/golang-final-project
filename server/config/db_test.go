package config_test

import (
	"fmt"
	"service/backup/databases/server/config"
	"testing"

	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err)
	}

}

func TestConnection(t *testing.T) {
	InitEnv()
	config.OpenDB()

}
