package cli

import (
	"os/exec"
	"sync"
	"time"
	"zschedule/cmd"
	"zschedule/configs"
	logger "zschedule/log"
)

type Scheduler struct {
	Aliases   []configs.CommandLineConfig
	Mutex     sync.Mutex
	WaitGroup sync.WaitGroup
}

func NewScheduler() (schedule *Scheduler, ere error) {
	var err error
	schedule = new(Scheduler)
	if schedule.Aliases, err = configs.ReadConfigFile(cmd.ConfigFile); err != nil {
		return nil, err
	}
	schedule.WaitGroup = sync.WaitGroup{}
	return
}

func (s *Scheduler) Run() {
	logger.WarningLogger.Println(len(s.Aliases))
	for _, schedule := range s.Aliases {
		s.WaitGroup.Add(1)
		go func() {
			defer s.WaitGroup.Done()
			switch schedule.Limit {
			case -1:
				var count int = 1
				for {
					logger.InfoLogger.Printf("[ %v ] command is going to start for the %d time", schedule.Name, count)
					time.Sleep(1 * time.Second)

					cmd := exec.Command("sh", "-c", schedule.Command)
					result, err := cmd.CombinedOutput()
					if err != nil {
						logger.ErrorLogger.Printf("error in run this command { %v }, error = %v|%s \n", schedule.Command, err, result)
					} else {
						logger.SuccessLogger.Printf("The output for this command { %v } is like :\n\n%s\n", schedule.Command, result)
					}

					logger.SuccessLogger.Printf("The [ %v ] command has been executed for the %d time", schedule.Name, count)
					count++
					time.Sleep(time.Second * time.Duration(schedule.Interval))
				}
			default:
				for count := range schedule.Limit {
					logger.InfoLogger.Printf("[ %v ] command is going to start for the %d/%d", schedule.Name, count+1, schedule.Limit)
					time.Sleep(1 * time.Second)

					cmd := exec.Command("sh", "-c", schedule.Command)
					result, err := cmd.CombinedOutput()
					if err != nil {
						logger.ErrorLogger.Printf("error in run this command { %v }, error = %v|%s \n", schedule.Command, err, result)
					} else {
						logger.SuccessLogger.Printf("The output for this command { %v } is like :\n\n%s\n", schedule.Command, result)
					}

					logger.SuccessLogger.Printf("The [ %v ] command has been executed for the %d/%d time", schedule.Name, count+1, schedule.Limit)
					time.Sleep(time.Second * time.Duration(schedule.Interval))
				}
			}
		}()
	}
	s.WaitGroup.Wait()
}
