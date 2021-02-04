// Copyright 2018-2019 Arthur Bailão. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE.md file.

package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	getter "github.com/hashicorp/go-getter"
)

// Download ...
func download(args *cliArgs) []string {
	filenames := make([]string, len(packages))
	for i, pkg := range packages {
		filename := path.Join(args.JarPath, pkg.Name())

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			downloadFile(filename, pkg.URL())
		}

		filenames[i] = filename
		log.Printf("%s downloaded", filename)
	}
	return filenames
}

func downloadFile(dst string, src string) {
	err := getter.GetFile(dst, src)
	if err != nil {
		panic(err)
	}
}

type mavenPackageInfo struct {
	Group    string
	Artifact string
	Version  string
}

func (pkg *mavenPackageInfo) URL() string {
	paths := strings.Split(pkg.Group, ".")
	paths = append(paths, pkg.Artifact, pkg.Version, pkg.Name())
	return "http://search.maven.org/remotecontent?filepath=" + strings.Join(paths, "/")
}

func (pkg *mavenPackageInfo) Name() string {
	return fmt.Sprintf("%s-%s.jar", pkg.Artifact, pkg.Version)
}

var packages = [...]mavenPackageInfo{
	{"commons-codec", "commons-codec", "1.9"},
	{"commons-logging", "commons-logging", "1.1.3"},
	{"commons-lang", "commons-lang", "2.6"},
	{"joda-time", "joda-time", "2.8.1"},
	{"com.amazonaws", "aws-java-sdk-core", "1.11.331"},
	{"com.amazonaws", "aws-java-sdk-cloudwatch", "1.11.331"},
	{"com.amazonaws", "aws-java-sdk-dynamodb", "1.11.331"},
	{"com.amazonaws", "aws-java-sdk-kinesis", "1.11.331"},
	{"com.amazonaws", "aws-java-sdk-kms", "1.11.331"},
	{"com.amazonaws", "aws-java-sdk-s3", "1.11.331"},
	{"com.amazonaws", "amazon-kinesis-client", "1.9.1"},
	{"com.fasterxml.jackson.core", "jackson-databind", "2.6.6"},
	{"com.fasterxml.jackson.core", "jackson-core", "2.6.6"},
	{"com.fasterxml.jackson.core", "jackson-annotations", "2.6.0"},
	{"com.fasterxml.jackson.dataformat", "jackson-dataformat-cbor", "2.6.6"},
	{"org.apache.httpcomponents", "httpclient", "4.5.2"},
	{"org.apache.httpcomponents", "httpcore", "4.4.4"},
	{"com.google.guava", "guava", "18.0"},
	{"com.google.protobuf", "protobuf-java", "2.6.1"},
}
