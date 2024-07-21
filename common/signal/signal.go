package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForExit() {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalCh
}
