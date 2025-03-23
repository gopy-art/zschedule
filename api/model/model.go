package model

import (
	"time"

	"gorm.io/gorm"
)

type ScheduleAPI struct {
	gorm.Model
	Name     string    `json:"name"`
	Command  string    `json:"command"`
	Interval int       `json:"interval"`
	Limit    int       `json:"limit"`
	Current  int       `json:"current"`
	Next     time.Time `json:"next"`
}

func (s *ScheduleAPI) Add(con *gorm.DB) error {
	return con.Model(&ScheduleAPI{}).Create(s).Error
}

func (s *ScheduleAPI) Delete(con *gorm.DB) error {
	return con.Model(&ScheduleAPI{}).Where("ID = ?", s.ID).Delete(&ScheduleAPI{}).Error
}

func (s *ScheduleAPI) SelectAll(con *gorm.DB) ([]ScheduleAPI, error) {
	var result []ScheduleAPI
	databaseResponse := con.Model(&ScheduleAPI{}).Find(&result)
	if databaseResponse.Error != nil {
		return nil, databaseResponse.Error
	}

	return result, nil
}

func (s *ScheduleAPI) Update(con *gorm.DB) error {
	return con.Model(&ScheduleAPI{}).Where("ID = ?", s.ID).Updates(s).Error
}
