package args

import (
	"duckhunt/pkg/file"
	"duckhunt/pkg/logger"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Init() will parse the command line arguments with the flag package
func Init() {
	// Parse the arguments
	var flagRegister bool
	flag.BoolVar(&flagRegister, "r", false, "")
	flag.Usage = func() {
		fmt.Printf(`duckhunt.exe [OPTIONS]

Options:
 -r Registers duckhunt as a scheduled task to start upon login
 `)
	}
	flag.Parse()

	// This was only made for Windows
	if runtime.GOOS != "windows" {
		logger.CatchCritErr(fmt.Errorf("not running under Windows"))
	}

	// This needs admin privs
	if !checkAdmin() {
		logger.CatchCritErr(fmt.Errorf("not running with administrator priviledges"))
	}

	// Register to startup if selected
	if flagRegister {

		// Admin priv check
		if !checkAdmin() {
			logger.CatchCritErr(fmt.Errorf("not running with administrator priviledges"))
		}

		// Get the path of the executable
		exePath, err := file.ExePath()
		if err != nil {
			logger.CatchCritErr(err)
		}

		// Command to register as a schedulded task
		cmd := exec.Command("schtasks", "/create", "/tn", "Start duckhunt", "/tr", exePath, "/sc", "onlogon")
		// Run
		err = cmd.Run()
		if err != nil {
			logger.CatchCritErr(fmt.Errorf("registering duckhunt as scheduled task failed with: %s", err))
		}

		logger.Write("Registered duckhunt as a scheduled task")
		os.Exit(0)
	}
}

// CheckAdmin() will return true if the current context is elevated
//
// https://blog.hadenes.io/post/how-to-request-admin-permissions-in-windows-through-uac-with-golang/
func checkAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}
