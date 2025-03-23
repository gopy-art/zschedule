package api

import (
	"os"
	"zschedule/api/database"
	"zschedule/api/handler"
	"zschedule/api/model"
	"zschedule/api/worker"
	logger "zschedule/log"

	"github.com/gofiber/fiber/v3"
	"github.com/gopy-art/zrediss/connection"
)

func Server() {
	// Initialize a new Fiber app
	app := fiber.New()
	db, err := database.NewDatabaseConnection()
	if err != nil {
		logger.ErrorLogger.Fatalf("error in make database instanse, error = %v", err)
	}
	if err := db.Connection(); err != nil {
		logger.ErrorLogger.Fatalf("error in connect to the postgres, error = %v \n", err)
	}
	db.CreateTables(&model.ScheduleAPI{})

	{
		wschedule := new(worker.ScheduleWorker)
		if os.Getenv("CACHE_ADDRESS") == "" {
			logger.ErrorLogger.Fatalf("CACHE_ADDRESS is empty in .env file")
		}
		wschedule.Cache = connection.RedisConnection{
			RedisAddress: os.Getenv("CACHE_ADDRESS"),
		}
		if err := wschedule.Cache.InitConnection(); err != nil {
			logger.ErrorLogger.Fatalf("error in connect to redis, error = %v", err)
		}
		go wschedule.Run(db.DB)
		go wschedule.CheckForRun(db.DB)
	}

	{
		handler := handler.SchedulerHandler{
			Database: db,
		}

		schedule := app.Group("/schedule")
		schedule.Post("/add", handler.AddHandler)
		schedule.Put("/update", handler.UpdateHandler)
		schedule.Get("/select", handler.SelectAllHandler)
		schedule.Delete("/delete", handler.DeleteHandler)
	}

	// Start the server on port -
	if os.Getenv("LISTEN_ADDRESS") == "" {
		logger.ErrorLogger.Fatalln("LISTEN_ADDRESS have to be set in the .env file")
	} else {
		if err := app.Listen(os.Getenv("LISTEN_ADDRESS")); err != nil {
			logger.ErrorLogger.Fatalf("error in run api server, error = %v \n", err)
		}
	}
}
