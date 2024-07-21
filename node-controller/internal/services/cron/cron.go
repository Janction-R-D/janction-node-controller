package cron

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	syncOnce = &sync.Once{}
	c        *crontab
)

func init() {
	c = &crontab{
		task: newCron(),
	}
}

func Run() error {
	syncOnce.Do(func() {
		c.start()
	})
	return c.err
}

func Stop() {
	c.stop()
}

type crontab struct {
	task *AppCron
	err  error
}

func (c *crontab) start() {
	err := c.task.start()
	if err != nil {
		c.err = errors.Wrap(err, "Failed to start cron job")
	}
	return
}

func (c *crontab) stop() {
	c.task.stop()
	logrus.Info("Stopped all accounts cron jobs")
}
