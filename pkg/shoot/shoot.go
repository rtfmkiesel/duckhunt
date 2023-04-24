package shoot

import (
	"duckhunt/pkg/config"
	"duckhunt/pkg/logger"
	"duckhunt/pkg/user32"
	"time"
)

var (
	lastAlert time.Time = time.Now()
)

// ShootDuck() will run when an attack is detected
// and alert/block if configured
func ShootDuck(cfg config.Cfg) {
	logger.LogWrite(cfg.AlertMsg)

	// Alert if selected and only alert max every 5s
	if cfg.Alert && lastAlert.Add(5*time.Second).Before(time.Now()) {
		lastAlert = time.Now()
		// Is a go func since the syscall
		// returns only if used clicks popup away
		go user32.MessageBox(cfg)
	}

	// Block the user inputs if selected
	err := user32.BlockInputFor(cfg)
	if err != nil {
		logger.CatchErr(err)
	}

	// For some reason Windows will have a
	// "half" ALT key set at this point, press to normalize
	user32.SendKey(user32.VK_ALT)
}
