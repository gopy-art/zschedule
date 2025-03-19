package main

import (
	"fmt"
	"os"
	"zschedule/cmd"
	logger "zschedule/log"
)

func init() {
	zMinioBaseDir, _ := os.Getwd()

	// set the flags
	cmd.Execute()

	// init logger
	if cmd.Logger == "stdout" {
		logger.InitLoggerStdout()
	} else if cmd.Logger == "file" {
		logger.InitLoggerFile(zMinioBaseDir + "/zminio.log")
	}
}

func main() {
	fmt.Println("hello world")
}
