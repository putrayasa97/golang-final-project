package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"service/backup/databases/client/model"
	"service/backup/databases/client/utils/dbbackup"
	"service/backup/databases/client/utils/logger"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Warn("Cannot load env file, using system env")
	}
}
func main() {
	InitEnv()

	var rootCmd = &cobra.Command{Use: "app"}

	var cmdBackupDBList = &cobra.Command{
		Use:   "backup-dblist",
		Short: "Show list databases : backup-dblist",
		Run: func(cmd *cobra.Command, args []string) {
			pathFile := model.PathFile{
				PathDBJson: "databases.json",
			}

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
		},
	}

	var cmdBackupDBName = &cobra.Command{
		Use:   "backup-dbname [name]",
		Short: "Run backup with name : backup-dbname [db_1,db_2]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dbbackup.BackupRunner(args[0])
		},
	}

	var cmdBackupScheduler = &cobra.Command{
		Use:   "backup-scheduler [time]",
		Short: "Run backup by scheduler : backup-scheduler [interval-minute]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			scheduler := cron.New()

			defer scheduler.Stop()

			logger.Info(fmt.Sprintln("Scheduler Run.."))
			scheduler.AddFunc("*/"+args[0]+" * * * *", func() { dbbackup.BackupRunner("") })

			go scheduler.Start()

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			<-sig
		},
	}

	rootCmd.AddCommand(cmdBackupDBList)
	rootCmd.AddCommand(cmdBackupDBName)
	rootCmd.AddCommand(cmdBackupScheduler)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
