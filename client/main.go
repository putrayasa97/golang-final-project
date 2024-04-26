package main

import (
	"fmt"
	"os"
	"os/signal"
	"service/backup/databases/client/utils/dbbackup"
	"service/backup/databases/client/utils/logger"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Warn("Cannot load env file, using system env")
	}
}
func main() {
	InitEnv()

	scheduler := cron.New()

	defer scheduler.Stop()

	logger.Info(fmt.Sprintln("Scheduler Run.."))
	scheduler.AddFunc("*/1 * * * *", func() { dbbackup.BackupRunner() })

	go scheduler.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
