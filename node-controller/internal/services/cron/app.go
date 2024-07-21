package cron

import (
	"github.com/go-co-op/gocron"
	"time"
)

type AppCron struct {
	scheduler *gocron.Scheduler
}

func newCron() *AppCron {
	return &AppCron{scheduler: gocron.NewScheduler(time.Local)}
}

func (c *AppCron) start() error {
	c.scheduler.StartAsync()
	return nil
}

func (c *AppCron) stop() {
	c.scheduler.Stop()
}
