package main

import (
	"duckhunt/pkg/args"
	"duckhunt/pkg/config"
	"duckhunt/pkg/logger"
	"duckhunt/pkg/shoot"
	"duckhunt/pkg/speed"
	"duckhunt/pkg/user32"
	"time"
)

func main() {

	// Parse args
	args.Init()
	// Load and parse the config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.CatchCritErr(err)
	}

	// Setup speed tracker
	i := 0
	history := speed.InitHistory(cfg.HistorySize, cfg.MaxInterval)
	timeOld := time.Now().UnixMilli()

	logger.LogWrite("duckhunt started")

	// Run forever
	for {

		// Reset i to loop around history
		if i >= cfg.HistorySize {
			i = 0
		}

		// Wait for a key press
		timePressed := user32.WaitForKey(cfg)

		// Save the interval
		history[i] = timePressed - timeOld

		// Advance the slice
		i = i + 1
		timeOld = timePressed

		// Trigger if the average intervals are smaller than specified
		if speed.CalcAvrgHistory(history) < cfg.MaxInterval {
			shoot.ShootDuck(cfg)
		}
	}
}
