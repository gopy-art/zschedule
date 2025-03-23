package worker

import (
	"fmt"
	"os/exec"
	"time"
	"zschedule/api/model"
	logger "zschedule/log"

	"github.com/gopy-art/zrediss/connection"
	"gorm.io/gorm"
)

var DataPool chan model.ScheduleAPI = make(chan model.ScheduleAPI, 1)

type ScheduleWorker struct {
	NumberOfWorkers int `json:"number_of_workers"`
	Cache           connection.RedisConnection
}

func (w *ScheduleWorker) CheckForRun(db *gorm.DB) {
	tmp := new(model.ScheduleAPI)
	for {
		schedules, err := tmp.SelectAll(db)
		if err != nil {
			logger.ErrorLogger.Printf("error in get list of schedules, error = %v \n", err)
		}
		for _, schedule := range schedules {
			if schedule.Next.Before(time.Now()) {
				if schedule.Limit == -1 {
					DataPool <- schedule
				} else if schedule.Limit != -1 && schedule.Current < schedule.Limit {
					DataPool <- schedule
				} else {
					continue
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (w *ScheduleWorker) Run(db *gorm.DB) {
	logger.InfoLogger.Printf("execute worker is up for handling schedules...")
	for task := range DataPool {
		if task.Next.After(time.Now()) {
			continue
		}

		logger.InfoLogger.Printf("[ %v ] command is going to start for the %d time", task.Name, task.Current+1)
		time.Sleep(1 * time.Second)

		cmd := exec.Command("sh", "-c", task.Command)
		result, err := cmd.CombinedOutput()
		if err != nil {
			logger.ErrorLogger.Printf("error in run this command { %v }, error = %v|%s \n", task.Command, err, result)
		} else {
			logger.SuccessLogger.Printf("The output for this command { %v } is like :\n\n%s\n", task.Command, result)
		}

		if _, err := w.Cache.SetKeyWithValue(
			fmt.Sprintf("%v - %v", task.Name, time.Now().Format("2006-01-02 15:04:05")),
			fmt.Sprintf("OUTPUT : %v", string(result)),
		); err != nil {
			logger.ErrorLogger.Printf("error in set output in the cache, error = %v\n", err)
		}

		logger.SuccessLogger.Printf("The [ %v ] command has been executed for the %d time", task.Name, task.Current+1)

		if err = db.Model(&model.ScheduleAPI{}).Where("ID = ?", task.ID).Updates(&model.ScheduleAPI{
			Current: task.Current + 1,
			Next:    time.Now().Add(time.Second * time.Duration(task.Interval)),
		}).Error; err != nil {
			logger.ErrorLogger.Printf("error in update the state of this schedule in database, error = %v\n", err)
		}
	}
}
