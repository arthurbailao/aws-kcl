package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type xmlMap map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// UnmarshalXML unmarshals the XML into a map of string to strings,
// creating a key in the map for each tag and setting it's value to the
// tags contents.
//
// The fact this function is on the pointer of Map is important, so that
// if m is nil it can be initialized, which is often the case if m is
// nested in another xml structurel. This is also why the first thing done
// on the first line is initialize it.
// source: https://go.dev/play/p/4Z2C-GF0E7
func (m *xmlMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = xmlMap{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[fmt.Sprintf("${%s}", e.XMLName.Local)] = e.Value
	}
	return nil
}

type Dependency struct {
	Group    string `xml:"groupId"`
	Artifact string `xml:"artifactId"`
	Version  string `xml:"version"`
}

func (d Dependency) URL(properties xmlMap) string {
	paths := strings.Split(d.Group, ".")
	paths = append(paths, d.Artifact, d.CommonVersion(properties), d.Name(properties))
	return "https://repo1.maven.org/maven2/" + strings.Join(paths, "/")
}

func (d Dependency) Name(properties xmlMap) string {
	return fmt.Sprintf("%s-%s.jar", d.Artifact, d.CommonVersion(properties))
}

func (d Dependency) CommonVersion(p xmlMap) string {
	if v, ok := p[d.Version]; ok {
		return v
	}
	return d.Version
}

type Project struct {
	XMLName      xml.Name     `xml:"project"`
	Properties   xmlMap       `xml:"properties"`
	Dependencies []Dependency `xml:"dependencies>dependency"`
}

// Download ...
func download(args *cliArgs) []string {
	data, err := os.ReadFile(args.PomPath)
	if err != nil {
		panic("failed to read pom.xml content")
	}

	p := Project{}
	err = xml.Unmarshal(data, &p)
	if err != nil {
		panic("failed to unmarshal pom.xml")
	}

	var filenames []string
	for _, d := range p.Dependencies {
		filename := path.Join(args.JarPath, d.Name(p.Properties))
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			downloadFile(filename, d.URL(p.Properties))
		}
		filenames = append(filenames, filename)
		log.Printf("%s downloaded", filename)
	}

	return filenames
}

func downloadFile(dst string, src string) {
	out, err := os.Create(dst)
	if err != nil {
		panic("failed to create jar file")
	}
	defer out.Close()

	resp, err := http.Get(src)
	if err != nil {
		panic("error requesting jar file")
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic("failed to download jar")
	}
}
