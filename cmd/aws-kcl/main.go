// Copyright 2018-2019 Arthur Bailão. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE.md file.

package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

const daemonClass = "com.amazonaws.services.kinesis.multilang.MultiLangDaemon"

func main() {
	args := parseArgs()
	jars := download(args)

	paths := append(jars, absDir(os.Args[0]), absDir(args.Properties))
	classpath := strings.Join(paths, string(os.PathListSeparator))

	var cmd = []string{
		args.Java,
		"-cp",
		classpath,
		daemonClass,
		args.Properties,
	}

	fmt.Println(strings.Join(cmd, " "))
	err := syscall.Exec(args.Java, cmd, os.Environ())
	if err != nil {
		panic(err)
	}
}
