package handler

import (
	"fmt"
	"zschedule/api/database"
	"zschedule/api/model"

	"github.com/gofiber/fiber/v3"
)

type SchedulerHandler struct {
	Database *database.DatabaseConfiguration
	Data     *model.ScheduleAPI
}

func (s *SchedulerHandler) AddHandler(c fiber.Ctx) error {
	s.Data = new(model.ScheduleAPI)

	if err := c.Bind().Body(s.Data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid request body",
		})
	}

	if s.Data.Command == "" || s.Data.Interval == 0 || s.Data.Limit == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "request body is out of format",
		})
	}

	if err := s.Data.Add(s.Database.DB); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": fmt.Sprintf("error in add schedule in database, error = %v", err),
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  200,
			"message": "schedule added successfully",
		})
	}
}

func (s *SchedulerHandler) UpdateHandler(c fiber.Ctx) error {
	s.Data = new(model.ScheduleAPI)

	if err := c.Bind().Body(s.Data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid request body",
		})
	}

	if s.Data.ID == uint(0) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "ID cannot be empty.",
		})
	}

	if err := s.Data.Update(s.Database.DB); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": fmt.Sprintf("error in update schedule in database, error = %v", err),
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  200,
			"message": "schedule updated successfully",
		})
	}
}

func (s *SchedulerHandler) DeleteHandler(c fiber.Ctx) error {
	s.Data = new(model.ScheduleAPI)

	if err := c.Bind().Body(s.Data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid request body",
		})
	}

	if s.Data.ID == uint(0) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "ID cannot be empty.",
		})
	}

	if err := s.Data.Delete(s.Database.DB); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": fmt.Sprintf("error in delete schedule in database, error = %v", err),
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  200,
			"message": "schedule deleted successfully",
		})
	}
}

func (s *SchedulerHandler) SelectAllHandler(c fiber.Ctx) error {
	if data, err := s.Data.SelectAll(s.Database.DB); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": fmt.Sprintf("error in select all schedules from database, error = %v", err),
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  200,
			"message": data,
		})
	}
}