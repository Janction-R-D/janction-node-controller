package runner

import (
	"common/buildinfo"
	"common/config"
	"flag"
	"os"
)

func MustParseStandardCommandArgs() *config.AppConf {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.json", "config path")
	printV := flag.Bool("v", false, "print version info")
	flag.Parse()
	if *printV {
		buildinfo.PrintBuildInfo()
		os.Exit(0)
	}

	conf, err := config.Read(configPath)
	if err != nil {
		panic(err)
	}

	return conf
}
