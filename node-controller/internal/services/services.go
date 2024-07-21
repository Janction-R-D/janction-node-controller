package services

import (
	"common/runner"
	"common/signal"
	"net"
	"node-controller/config"
	"node-controller/internal/routers"
	"node-controller/internal/services/cron"
	"strconv"
	"sync"

	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

type App struct{}

func (a *App) Run(runner *runner.Runner) (err error) {
	config.MustInitConfig(runner.Conf)
	Init()
	err = cron.Run()
	if err != nil {
		log.WithError(err).Fatal("Failed to run cron jobs")
	}
	var addr = net.JoinHostPort(config.AppHost, strconv.Itoa(config.AppPort))
	err = newApp().Run(iris.Addr(addr))
	if err != nil {
		log.WithError(err).Fatal("failed to run app")
	}
	signal.WaitForExit()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cron.Stop()
	}()
	return nil
}

func newApp() *iris.Application {
	var app = iris.New()
	routers.Register(app)
	return app
}
