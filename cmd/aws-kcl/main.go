package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

const daemonClass = "com.amazonaws.services.kinesis.multilang.MultiLangDaemon"

func main() {
	args := ParseArgs()
	jars := Download(args)

	paths := append(jars, AbsDir(os.Args[0]), AbsDir(args.Properties))
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
