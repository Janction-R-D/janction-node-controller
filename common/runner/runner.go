package runner

import (
	"common/buildinfo"
	"common/config"
	"common/logfmt"
	"common/pprof"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Runner struct {
	Name string
	Conf *config.AppConf
}

type Runnable interface {
	Run(runner *Runner) error
}

type LogSetter interface {
	LogOutput(conf *config.AppConf) io.Writer
}

type LogFormatter interface {
	LogFormat() logrus.Formatter
}

func New(name string, conf *config.AppConf) *Runner {
	return &Runner{
		Name: name,
		Conf: conf,
	}
}

func (r *Runner) adjustLoglevelTask() {
	_, name := filepath.Split(r.Name)
	fileName := filepath.Join(os.TempDir(), name+"_log_level")
	defaultLevel := logrus.GetLevel()

	logrus.WithFields(logrus.Fields{"level": defaultLevel, "level_file": fileName}).Info("You can change log level in this file")

	go func() {
		tick := time.NewTicker(30 * time.Second)
		for range tick.C {
			b, err := os.ReadFile(fileName)
			if err != nil {
				changeLevel(defaultLevel)
				continue
			}
			fileCnt := strings.TrimSpace(string(b))
			level, err := logrus.ParseLevel(fileCnt)
			if err != nil {
				level = defaultLevel
			}
			changeLevel(level)
		}
	}()
}

func changeLevel(level logrus.Level) {
	if logrus.GetLevel() != level {
		logrus.WithFields(logrus.Fields{
			"from": logrus.GetLevel(),
			"to":   level,
		}).Info("log level is changed")
	}
	logrus.SetLevel(level)
}

func (r *Runner) setupLog(conf *config.AppConf, runnable Runnable) {
	logFormat, ok := runnable.(LogFormatter)
	if ok {
		logrus.SetFormatter(logFormat.LogFormat())
	} else {
		log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
		logrus.SetFormatter(&logfmt.BaFormatter{})
		logrus.SetReportCaller(true)
	}

	logOutput, ok := runnable.(LogSetter)
	if ok {
		logrus.SetOutput(logOutput.LogOutput(conf))
	} else {
		logBaseDir := conf.GetString("log.log_base_dir", "/home/logs")
		_ = os.MkdirAll(logBaseDir, os.ModePerm)

		logMaxSizeMB := conf.GetInt("log.log_max_size_mb", 200)
		logMaxBackup := conf.GetInt("log.log_max_backup", 5)
		logCompress := conf.GetBool("log.compress", true)

		logPath := filepath.Join(logBaseDir, filepath.Base(r.Name)+".log")
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    int(logMaxSizeMB), // in MB
			MaxBackups: int(logMaxBackup),
			Compress:   logCompress,
		})
	}

	logLevel := conf.GetString("log.level", "info")
	levelStr := strings.ToLower(logLevel)
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		level = logrus.WarnLevel
	}
	logrus.SetLevel(level)
}

func (r *Runner) serveProfile(conf *config.AppConf) {
	_, name := filepath.Split(r.Name)
	addrKey := name + "_addr"
	profileAddr := conf.GetString("profile."+addrKey, "")
	if profileAddr == "" {
		return
	}
	go func() {
		logrus.Info("pprof server is listen at ", profileAddr)
		profSrv := pprof.NewServer()
		logrus.WithError(profSrv.ListenAndServe(profileAddr)).Error("profSrv exited")
	}()
}

func (r *Runner) Run(runnable Runnable) error {
	rand.Seed(time.Now().UnixNano())

	r.setupLog(r.Conf, runnable)

	logrus.WithField("pid", os.Getpid()).Infof("%s [start]", filepath.Base(r.Name))
	defer func() {
		logrus.WithField("pid", os.Getpid()).Infof("%s [quit]", filepath.Base(r.Name))
	}()

	r.adjustLoglevelTask()
	r.serveProfile(r.Conf)

	buildinfo.LogBuildInfo()

	return runnable.Run(r)
}
