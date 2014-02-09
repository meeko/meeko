// Copyright (c) 2013 The ciderapp AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"github.com/cihub/seelog"
	"github.com/tchap/gocli"
	"os"
)

var ciderapp = gocli.NewApp("ciderapp")

func init() {
	ciderapp.UsageLine = "ciderapp [-debug] [-endpoint ENDPOINT] SUBCMD"
	ciderapp.Version = "0.0.1"
	ciderapp.Short = "Cider applications management utility"
	ciderapp.Long = `
  ciderapp is a command line utility for managing local Cider instance,
  or rather the Cider applications running on it.

  This tool expects the local Cider instance's management token to be saved
  in .cider_token file places in the current user's home directory.`
	ciderapp.Flags.BoolVar(&fdebug, "debug", fdebug, "print debug output")
	ciderapp.Flags.StringVar(&fendpoint, "endpoint", fendpoint, "Cider ZeroMQ 3.x RPC endpoint")
}

var (
	fdebug    bool
	fendpoint string
)

func main() {
	seelog.ReplaceLogger(seelog.Disabled)
	ciderapp.Run(os.Args[1:])
}
