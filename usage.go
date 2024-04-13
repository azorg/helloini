// File: "usage.go"

package main

import (
	"fmt"
	"os"
)

// Print usage and exit
func printUsage() {
	fmt.Println(APP_DESCRIPTION + " v" + Version + `
Usage: ` + APP_NAME + ` [flags] command [options]

Flags:
  -c <config_file> - Path to config INI file

Logger flags:
  -log  <level> - Log level (trace/debug/info/notify/warn/error/fatal)
  -slog         - Use structured text logger (slog)
  -jlog         - Use structured JSON logger (slog)
  -tlog         - Use tinted (colorized) logger (tint)
  -lsrc         - Force log source file name and line number
  -lpkg         - Force log source directory/file name and line number
  -ltime        - Force add time to log
  -ltimefmt     - Override log time format (e.g. 15:04:05.999)

Local commands:
  -h|--h               - Show short help about options only and exit
  help|-help|--help    - Show full help and exit
  version|-v|--version - Show version and exit
  mkconf [conf]        - Generate default config INI file and exit
  showconf             - Show configuration as INI file to stdout and exit

Commands:
  start - start application (by default)

Enviroment variable:
  ` + CONF_ENV + ` - path to default config INI file
`)
	os.Exit(0)
}

// Print version and exit
func printVersion() {
	fmt.Println(APP_NAME + " version " + Version)
	os.Exit(0)
}

// EOF: "usage.go"
