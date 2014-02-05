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
	"fmt"
	"io"
	"os"
	"strings"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/wsxiaoys/terminal/color"
)

var (
	OK   = color.Sprintf("[ @{g}OK@{|} ]\n")
	SKIP = color.Sprintf("[ @{y}SKIP@{|} ]\n")
	FAIL = color.Sprintf("[ @{!}@{r}FAIL@{|} ]\n")
)

var replyCache = make(map[string]string)

func promptUser(prompt string, secret bool) (value string, err error) {
	if val, ok := replyCache[prompt]; ok {
		return val, nil
	}

	defer func() {
		if err == nil {
			replyCache[prompt] = value
		}
	}()

	file, err := os.OpenFile("/dev/tty", os.O_RDWR, 0600)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(prompt)
	if err != nil {
		return
	}

	var bs []byte
	if secret {
		bs, err = terminal.ReadPassword(int(file.Fd()))
		if err != nil {
			return "", err
		}

		_, err = file.Write([]byte("\n"))
		if err != nil {
			return "", err
		}
	} else {
		bs = make([]byte, 80)
		n, err := io.LimitReader(file, 80).Read(bs)
		if err != nil {
			return "", err
		}
		bs = bs[:n]
	}

	return strings.TrimSuffix(string(bs), "\n"), nil
}

func debug(a ...interface{}) {
	if fdebug {
		fmt.Print(a...)
	}
}
