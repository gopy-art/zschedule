package api

import (
	logger "zschedule/log"

	"github.com/gofiber/fiber/v3"
)

func Server() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World!")
	})

	// Start the server on port -
	err := app.Listen(":3030")	
	if err != nil {
		logger.ErrorLogger.Fatalf("error in run api server, error = %v \n", err)
	}
}
