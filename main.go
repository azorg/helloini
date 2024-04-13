// File: "main.go"

package main

import (
	//"fmt"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/azorg/xlog"
)

func main() {
	// Create default config structure
	conf := NewConf()

	// Setup logger by default
	xlog.Setup(conf.Log)

	// Ceate command line options
	opt := NewOpt()

	// Parse command line options
	opt.Parse()

	// Create default config and exit if `mkconf` command set
	//conf.SaveAndExit(opt.Conf)

	// Read config ini file
	_ = conf.Load(opt.Conf)

	// Add command line options to config structure
	conf.AddOpt(opt)

	// Setup logger
	xlog.Setup(conf.Log)

	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
		xlog.Error("os.Getwd()", "err", err)
	}

	saveConf := confFileToSave()
	xlog.Info(
		APP_NAME+" start",
		"version", Version,
		"logLvl", xlog.GetLvl(),
		"findConf", opt.Conf,
		"saveConf", saveConf,
		"pwd", pwd)

	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()

	// Save configuration/state (update INI file)
	err = conf.Update(saveConf)
	if err != nil {
		xlog.Error("can't save INI file", "saveConf", saveConf)
	} else {
		xlog.Info("updated INI file", "saveConf", saveConf)
	}
}

// EOF: "main.go"
