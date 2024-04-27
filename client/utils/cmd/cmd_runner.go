package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"service/backup/databases/client/model"
	"service/backup/databases/client/utils/dbbackup"
	"service/backup/databases/client/utils/dbconfig"
	"service/backup/databases/client/utils/logger"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func CommandLine() {

	var rootCmd = &cobra.Command{Use: "app"}

	var cmdBackupDB = &cobra.Command{
		Use:   "backup-db",
		Short: "Run backup : backup-db",
		Run: func(cmd *cobra.Command, args []string) {
			dbbackup.BackupRunner("")
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
		Use:   "backup-dbscheduler [time]",
		Short: "Run backup by scheduler : backup-dbscheduler [interval-minute]",
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

	pathFile := model.PathFile{
		PathDBJson: "databases.json",
	}

	var cmdBackupDBList = &cobra.Command{
		Use:   "backup-dblist",
		Short: "Show config databases : backup-dblist",
		Run: func(cmd *cobra.Command, args []string) {
			dbconfig.ListDatabaseConfig(pathFile)
		},
	}

	var cmdBackupDBAdd = &cobra.Command{
		Use:   "backup-dbadd [db_name] [db_host] [db_port] [db_username] [db_password]",
		Short: "Add config databases : backup-dbadd [db_name] [db_host] [db_port] [db_username] [db_password]",
		Args:  cobra.ExactArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			dbconfig.AddDatabaseConfig(args, pathFile)
		},
	}

	var cmdBackupDBDelete = &cobra.Command{
		Use:   "backup-dbdelete [db_name]",
		Short: "Delete config databases : backup-dbdelete [db_name]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dbconfig.DeleteDatabaseConfig(args[0], pathFile)
		},
	}

	rootCmd.AddCommand(cmdBackupDBAdd)
	rootCmd.AddCommand(cmdBackupDBDelete)
	rootCmd.AddCommand(cmdBackupDBList)
	rootCmd.AddCommand(cmdBackupDB)
	rootCmd.AddCommand(cmdBackupDBName)
	rootCmd.AddCommand(cmdBackupScheduler)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
