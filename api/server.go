package api

import (
	"zschedule/api/database"
	"zschedule/api/handler"
	"zschedule/api/model"
	logger "zschedule/log"

	"github.com/gofiber/fiber/v3"
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
	if err := app.Listen(":3030"); err != nil {
		logger.ErrorLogger.Fatalf("error in run api server, error = %v \n", err)
	}
}
