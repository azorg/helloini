// File: "const.go"

package main

const (
	APP_NAME        = "helloini"
	APP_DESCRIPTION = "hello world application with INI configuration"
	VERSION_MAJOR   = "x"
	VERSION_MINOR   = "y"
	VERSION_BUILD   = "z"
)

var Version = VERSION_MAJOR + "." + VERSION_MINOR + "." + VERSION_BUILD
var Hash string

func init() {
	if Hash != "" {
		Version += "-" + Hash
	}
}

// EOF: "const.go"
