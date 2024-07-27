package ticker

import (
	"janction/logic"
	"time"

	"go.uber.org/zap"
)

func FetchJobTicker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := logic.FetchJob()
			if err != nil {
				zap.L().Error("Ticker report error happend", zap.Error(err))
				continue
			}
		}
	}
}
