// File: "conf.go"

package main

import (
	"flag"
	"os"
	"path"
	"strconv"

	"github.com/azorg/xlog"
	"gopkg.in/ini.v1"
)

// Defaults
const (
	CONF_FILE     = "hello.ini"     // config file name
	CONF_ENV      = "HELLO_INI"     // enviroment variable of path to config
	HOME_ENV      = "HOME"          // enviroment variable of home directory
	LOG_PREFIX    = "LOG_"          // enviroment logger prefix
	CONF_DIR_PWD  = "."             // curent directory
	CONF_DIR_HOME = ".config/hello" // config directory into home directory
	CONF_DIR_ETC  = "/etc/hello"    // config directory into /etc

	DEF_W = 300 // default witdh
	DEF_H = 200 // default height
)

// Command line options
type Opt struct {
	Conf string    // path to config INI file
	Log  *xlog.Opt // logger options
}

// Configuration structure
type Conf struct {
	// Position:
	X int `ini:"x"`
	Y int `ini:"y"`

	// Size:
	W int `ini:"w"`
	H int `ini:"h"`

	// Logger settings
	Log xlog.Conf
}

// Return directories list for find config file
func confDirs() []string {
	dirs := []string{}
	if CONF_DIR_PWD != "" {
		dirs = append(dirs, CONF_DIR_PWD)
	}

	if home := os.Getenv(HOME_ENV); home != "" && home != "." {
		dirs = append(dirs, path.Join(home, CONF_DIR_HOME))
	}

	if CONF_DIR_ETC != "" {
		dirs = append(dirs, CONF_DIR_ETC)
	}

	return dirs
}

// Find config file to load
func confFile() string {
	for _, dir := range confDirs() {
		conf := path.Join(dir, CONF_FILE)
		f, err := os.Open(conf)
		if err != nil {
			continue
		}
		fi, err := f.Stat()
		if err != nil {
			continue
		}
		if fi.IsDir() {
			continue
		}
		xlog.Debug("select config file", "conf", conf)
		return conf
	} // for
	return ""
}

// Find config file to save
func confFileToSave() string {
	if home := os.Getenv(HOME_ENV); home != "" && home != "." {
		dir := path.Join(home, CONF_DIR_HOME)
		err := os.MkdirAll(dir, 0750)
		if err != nil {
			xlog.Crit("can't create config directory", "dir", dir, "err", err)
			return CONF_FILE
		}
		return path.Join(dir, CONF_FILE)
	}
	return CONF_FILE
}

// Create command line options structure
func NewOpt() *Opt {
	o := &Opt{}

	if v := os.Getenv(CONF_ENV); v != "" {
		o.Conf = v
		xlog.Debug("select path to config file from enviroment", CONF_ENV, v)
	}

	flag.StringVar(&o.Conf, "c", o.Conf, "path to config INI file")

	o.Log = xlog.NewOpt()
	return o
}

// Create default configuration structure
func NewConf() *Conf {
	c := &Conf{
		X:   0,
		Y:   0,
		W:   DEF_W,
		H:   DEF_H,
		Log: xlog.NewConf(),
	}

	// Setup logger by default
	c.Log.Level = "debug"
	c.Log.Tint = true
	//c.Log.Src = true
	//c.Log.TimeTint = "timeOnly"

	return c
}

// Parse command line options
func (o *Opt) Parse() {
	for _, opt := range os.Args[1:] {
		if opt == "-help" || opt == "--help" || opt == "help" {
			printUsage()
		} else if opt == "-v" || opt == "--version" || opt == "version" {
			printVersion()
		} else if opt[0:1] != "-" {
			break // abort by first command
		}
	}

	flag.Parse()

	// Find config file if it not set
	if o.Conf == "" {
		o.Conf = confFile()
	}
}

// Add options to config structure
func (c *Conf) AddOpt(o *Opt) {
	// Get logger settings from enviroment
	xlog.Env(&c.Log, LOG_PREFIX)

	// Add logger options to configuration
	xlog.AddOpt(o.Log, &c.Log)
}

// Load configuration from INI file
func (c *Conf) Load(fileName string) error {
	if fileName == "" {
		return nil // do nothing
	}

	// Load INI file
	cfg, err := ini.Load(fileName)
	if err != nil {
		xlog.Error("fail to read INI file", "fileName", fileName, "err", err)
		return err
	}

	s := cfg.Section("position")
	c.X = s.Key("x").MustInt(c.X)
	c.Y = s.Key("y").MustInt(c.Y)

	s = cfg.Section("size")
	c.W = s.Key("w").MustInt(c.W)
	c.H = s.Key("h").MustInt(c.H)

	s = cfg.Section("log")
	l := &c.Log
	c.Log = xlog.Conf{
		File:     s.Key("file").MustString(l.File),
		FileMode: s.Key("file-mode").MustString(l.FileMode),
		Level: s.Key("level").In(l.Level, []string{
			"flood", "trace", "debug", "info", "notice", "warn",
			"error", "crit", "fatal", "panic", "silent"}),
		Slog:     s.Key("slog").MustBool(l.Slog),
		JSON:     s.Key("json").MustBool(l.JSON),
		Tint:     s.Key("tint").MustBool(l.Tint),
		Time:     s.Key("time").MustBool(l.Time),
		TimeUS:   s.Key("time-us").MustBool(l.TimeUS),
		TimeTint: s.Key("time-tint").MustString(l.TimeTint),
		Src:      s.Key("src").MustBool(l.Src),
		SrcLong:  s.Key("src-long").MustBool(l.SrcLong),
		NoLevel:  s.Key("no-level").MustBool(l.NoLevel),
		NoColor:  s.Key("no-color").MustBool(l.NoColor),
		Prefix:   s.Key("prefix").MustString(l.Prefix),
		AddKey:   s.Key("add-key").MustString(l.AddKey),
		AddValue: s.Key("add-value").MustString(l.AddValue),
	}

	return nil
}

// Save configuration to INI file
func (c *Conf) Update(fileName string) error {
	cfg, err := ini.Load(fileName)
	if err != nil {
		xlog.Notice("fail to read INI file before update",
			"fileName", fileName, "err", err)
		cfg = ini.Empty()
	}

	s := cfg.Section("position")
	s.Key("x").SetValue(strconv.Itoa(c.X))
	s.Key("y").SetValue(strconv.Itoa(c.Y))

	s = cfg.Section("size")
	s.Key("w").SetValue(strconv.Itoa(c.W))
	s.Key("h").SetValue(strconv.Itoa(c.H))

	//s = cfg.Section("log")
	//s.Key("file").SetValue

	err = cfg.SaveTo(fileName)
	if err != nil {
		xlog.Error("fail to save INI file", "fileName", fileName, "err", err)
	}
	return err
}

// EOF: "conf.go"
