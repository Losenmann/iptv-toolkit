package scheduler

import (
    "fmt"
//    "time"
    "log/slog"
    "github.com/go-co-op/gocron/v2"
    "iptv-toolkit/main/setup"
)

func Main(expression string) {
	if s, err := gocron.NewScheduler(); err != nil {
        if *setup.LogLVL <= 2 {
            slog.Warn(fmt.Sprintf("%v", err))
        }
    } else {
	    _, err := s.NewJob(
		    gocron.CronJob(expression, false),
		    gocron.NewTask(func() {Task()}),
	    )
	    if err != nil {
		    if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
	    }
	    //j.ID()
	    s.Start()
        if *setup.LogLVL <= 1 {
            slog.Info("A regular task has been created. Scheduler Expression: " + expression)
        }

	// block until you are ready to shut down
//	    select {
//	    case <-time.After(5*time.Minute):
//	    }

	// when you're done, shut it down
    /*
	    err = s.Shutdown()
	    if err != nil {
            if *setup.LogLVL <= 2 {
                slog.Warn(fmt.Sprintf("%v", err))
            }
	    }
    */
    }
}