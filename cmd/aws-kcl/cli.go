package main

import (
	"flag"
	"os"
	"path"
	"path/filepath"
)

type cliArgs struct {
	Properties string
	Java       string
	JarPath    string
}

const (
	propertiesUsage = "properties file with multi-language daemon options"
	javaUsage       = "path to java executable - defaults to using JAVA_HOME environment variable to get java path"
	jarPathUsage    = "path where all multi-language daemon jar files will be downloaded (optional)"
)

func parseArgs() *cliArgs {
	args := new(cliArgs)
	flag.StringVar(&args.Properties, "properties", "", propertiesUsage)
	flag.StringVar(&args.Java, "java", javaPath(), javaUsage)
	flag.StringVar(&args.JarPath, "jar_path", "", jarPathUsage)
	flag.Parse()

	if len(args.Properties) == 0 || len(args.Java) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(args.JarPath) == 0 {
		args.JarPath = jarPath()
	}

	return args
}

func absDir(path string) string {
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		panic(err)
	}
	return dir
}

func javaPath() string {
	javaHome, ok := os.LookupEnv("JAVA_HOME")
	if ok {
		return path.Join(javaHome, "bin", "java")
	}
	return ""
}

func jarPath() string {
	p := path.Join(absDir(os.Args[0]), "lib", "jar")
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return p
}
