// Copyright (c) 2013 The mk AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

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
	if flagDebug {
		fmt.Print(a...)
	}
}
