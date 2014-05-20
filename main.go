// Copyright (c) 2013 The meeko AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/cihub/seelog"
	"github.com/tchap/gocli"
)

var app = gocli.NewApp("meeko")

func init() {
	if value := os.Getenv("MEEKO_CONFIG"); value != "" {
		flagConfig = value
	}

	app.UsageLine = "meeko [-debug] [-config=PATH] SUBCMD"
	app.Version = "0.0.1"
	app.Short = "Meeko agents management utility"
	app.Long = `
  meeko is a command line utility for managing Meeko agents.
  Check the list of subcommands to see what actions are available.

  meeko expects a configuration file called .meekorc to be present
  in the user's home directory. The path can be ovewritten with -config.

ENVIRONMENT:
  MEEKO_CONFIG - can be used instead of -config
`

	app.Flags.BoolVar(&flagDebug, "debug", flagDebug, "print debug output")
	app.Flags.StringVar(&flagConfig, "config", flagConfig, "configuration file path")
}

var (
	flagDebug  bool
	flagConfig string = MustDefaultConfigPath()
)

func main() {
	seelog.ReplaceLogger(seelog.Disabled)
	app.Run(os.Args[1:])
}
