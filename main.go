package main

import (
	"fmt"
	"janction/dao/postgres"
	"janction/logger"
	"janction/pkg/snowflake"
	"janction/router"
	"janction/setting"
	"janction/ticker"

	"go.uber.org/zap"
)

func main() {
	// init setting
	if err := setting.Init(); err != nil {
		fmt.Println("Init setting failed, ", err)
		return
	}

	// init logger
	if err := logger.Init(setting.Config.LogConfig, setting.Config.Mode); err != nil {
		fmt.Println("Init logger failed, ", err)
		return
	}
	defer zap.L().Sync()

	// init snowflake
	if err := snowflake.Init(setting.Config.StartTime, setting.Config.MachineID); err != nil {
		fmt.Println("Init snowflake failed, ", err)
	}

	// init postgres connection
	if err := postgres.Init(setting.Config.PostgresConfig); err != nil {
		fmt.Println("Init postgres failed, ", err)
		return
	}
	defer postgres.Close()

	go ticker.FetchJobTicker()

	r := router.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", setting.Config.Port)); err != nil {
		fmt.Println("Run server failed, ", err)
		return
	}
}
