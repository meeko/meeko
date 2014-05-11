// Copyright (c) 2013 The mk AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"github.com/cihub/seelog"
	"github.com/tchap/gocli"
	"os"
)

var mk = gocli.NewApp("mk")

func init() {
	mk.UsageLine = "mk [-debug] [-endpoint ENDPOINT] SUBCMD"
	mk.Version = "0.0.1"
	mk.Short = "Cider applications management utility"
	mk.Long = `
  mk is a command line utility for managing local Cider instance,
  or rather the Cider applications running on it.

  This tool expects the local Cider instance's management token to be saved
  in .cider_token file places in the current user's home directory.`
	mk.Flags.BoolVar(&fdebug, "debug", fdebug, "print debug output")
	mk.Flags.StringVar(&fendpoint, "endpoint", fendpoint, "Cider ZeroMQ 3.x RPC endpoint")
}

var (
	fdebug    bool
	fendpoint string
)

func main() {
	seelog.ReplaceLogger(seelog.Disabled)
	mk.Run(os.Args[1:])
}
