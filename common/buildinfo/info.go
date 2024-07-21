package buildinfo

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	Version   string
	BuildTime string
)

func PrintBuildInfo() {
	fmt.Println("version:", Version)
	fmt.Println("build_time:", BuildTime)
}

func LogBuildInfo() {
	logrus.WithFields(logrus.Fields{
		"version":    Version,
		"build_time": BuildTime,
	}).Info("build info")
}
