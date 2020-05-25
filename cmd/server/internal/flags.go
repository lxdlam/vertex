package internal

import (
	"flag"
)

var (
	Verbose  bool
	Debug    bool
	ConfPath string
)

func init() {
	flag.BoolVar(&Verbose, "v", false, "print verbose log")
	flag.BoolVar(&Debug, "debug", false, "enter debug mode")
	flag.StringVar(&ConfPath, "conf_path", "./config.toml", "config file path")
}

func ParseFlags() {
	flag.Parse()
}
