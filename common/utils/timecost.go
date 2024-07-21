package utils

import (
	"fmt"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func TimeCost() func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		pc, _, _, _ := runtime.Caller(1)
		funcObj := runtime.FuncForPC(pc)
		fmt.Println(fmt.Sprintf("%60s took %.6fms", funcObj.Name(), float64(elapsed)/(1000*1000)))
	}
}

func TimeCostWithLogrus() func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		pc, _, _, _ := runtime.Caller(1)
		funcObj := runtime.FuncForPC(pc)
		logrus.Debug(fmt.Sprintf("%60s took %.6fms", funcObj.Name(), float64(elapsed)/(1000*1000)))
	}
}

func TimeCostWithLoggerAndText(logger *logrus.Logger, text string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		pc, _, _, _ := runtime.Caller(1)
		funcObj := runtime.FuncForPC(pc)
		logger.Debug(fmt.Sprintf("%60s took %.6fms(%s)", funcObj.Name(), float64(elapsed)/(1000*1000), text))
	}
}
