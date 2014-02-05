/*
   The MIT License (MIT)

   Copyright (c) 2013 Ondřej Kupka

   Permission is hereby granted, free of charge, to any person obtaining a copy of
   this software and associated documentation files (the "Software"), to deal in
   the Software without restriction, including without limitation the rights to
   use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
   the Software, and to permit persons to whom the Software is furnished to do so,
   subject to the following conditions:

   The above copyright notice and this permission notice shall be included in all
   copies or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
   FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
   COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
   IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
   CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

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
