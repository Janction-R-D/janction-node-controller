package main

import (
	"node-controller/internal/services"
	"common/runner"
	"os"
)

func main() {
	var appName = os.Args[0]
	var appConf = runner.MustParseStandardCommandArgs()
	var err = runner.New(appName, appConf).Run(&services.App{})
	if err != nil {
		panic(err)
	}
}
