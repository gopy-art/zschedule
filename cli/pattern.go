package cli

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
	"zschedule/cmd"
	"zschedule/configs"
	logger "zschedule/log"

	"github.com/gopy-art/zrediss/connection"
)

type Scheduler struct {
	Aliases   []configs.CommandLineConfig
	cache     connection.RedisConnection
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

	if os.Getenv("CACHE_ADDRESS") == "" {
		return nil, fmt.Errorf("CACHE_ADDRESS is empty in .env file")
	}
	schedule.cache = connection.RedisConnection{
		RedisAddress: os.Getenv("CACHE_ADDRESS"),
	}
	if err := schedule.cache.InitConnection(); err != nil {
		logger.ErrorLogger.Fatalf("error in connect to redis, error = %v", err)
	}
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

					if _, err := s.cache.SetKeyWithValue(
						fmt.Sprintf("%v - %v", schedule.Name, time.Now().Format("2006-01-02 15:04:05")),
						fmt.Sprintf("OUTPUT : %v", string(result)),
					); err != nil {
						logger.ErrorLogger.Printf("error in set output in the cache, error = %v\n", err)
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

					if _, err := s.cache.SetKeyWithValue(
						fmt.Sprintf("%v - %v", schedule.Name, time.Now().Format("2006-01-02 15:04:05")),
						fmt.Sprintf("OUTPUT : %v", string(result)),
					); err != nil {
						logger.ErrorLogger.Printf("error in set output in the cache, error = %v\n", err)
					}

					logger.SuccessLogger.Printf("The [ %v ] command has been executed for the %d/%d time", schedule.Name, count+1, schedule.Limit)
					time.Sleep(time.Second * time.Duration(schedule.Interval))
				}
			}
		}()
	}
	s.WaitGroup.Wait()
}
