package main

import (
	"fmt"
	"os"
	"zschedule/api"
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
	if cmd.Type == "api" {
		if cmd.EnvFile != "" {
			if err := godotenv.Load(cmd.EnvFile); err != nil {
				fmt.Printf("error in load the .env file, error = %v\n", err)
				os.Exit(1)
			}
		}
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

	} else if cmd.Type == "api" {
		api.Server()
	}
}
