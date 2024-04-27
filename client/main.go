package main

import (
	"service/backup/databases/client/utils/cmd"
	"service/backup/databases/client/utils/logger"

	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Warn("Cannot load env file, using system env")
	}
}
func main() {
	InitEnv()
	cmd.CommandLine()
}
