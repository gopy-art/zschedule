package main

import (
	"fmt"
	"os"
	"zschedule/api"
	"zschedule/cli"
	"zschedule/cmd"
	logger "zschedule/log"

	"github.com/joho/godotenv"
)

func init() {
	zMinioBaseDir, _ := os.Getwd()

	// set the flags
	cmd.Execute()

	// validate flags
	cmd.Validate()

	// init env
	if cmd.EnvFile != "" {
		if err := godotenv.Load(cmd.EnvFile); err != nil {
			fmt.Printf("error in load the .env file, error = %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("--env should be set")
		os.Exit(1)
	}

	// init logger
	if cmd.Logger == "stdout" {
		logger.InitLoggerStdout()
	} else if cmd.Logger == "file" {
		logger.InitLoggerFile(zMinioBaseDir + "/zminio.log")
	}
}

func main() {
	if cmd.Type == "cli" {
		schedulerConfig, err := cli.NewScheduler()
		if err != nil {
			logger.ErrorLogger.Fatalf("error in make new instance of the scheduler, error = %v \n", err)
		}
		schedulerConfig.Run()
	} else if cmd.Type == "api" {
		api.Server()
	}
}
