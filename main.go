package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Children struct {
	BinderItem []BinderItem `xml:"BinderItem"`
}

type BinderItem struct {
	UUID     string   `xml:"UUID,attr"`
	Title    string   `xml:"Title"`
	Children Children `xml:"Children"`
}

type ScrivenerProject struct {
	XMLName xml.Name `xml:"ScrivenerProject"`
	Binder  struct {
		Text       string `xml:",chardata"`
		BinderItem []BinderItem
	} `xml:"Binder"`
}

var paths = map[string]string{}

func main() {
	changedFiles := []string{}

	args := append([]string{"diff"}, os.Args[1:]...)

	matches, err := filepath.Glob("*.scrivx")
	if err != nil {
		log.Fatal(err)
	}
	if len(matches) == 0 {
		log.Fatalln("scrivgit must be run inside the .scriv directory")
	}
	if len(matches) > 1 {
		log.Fatalln("way too many .scrivx files for me to understand")
	}

	scrivx := matches[0]
	f, err := os.Open(scrivx)
	if err != nil {
		log.Fatalf("Failed to open %s: %v\n", scrivx, err)
	}

	decoder := xml.NewDecoder(f)

	var toc ScrivenerProject

	err = decoder.Decode(&toc)
	if err != nil {
		log.Fatal(err)
	}

	walkToc([]string{}, toc.Binder.BinderItem)

	nameargs := append(args, "--name-only", "--color=never", "--relative", ".")
	cmd := exec.Command("git", nameargs...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		changedFiles = append(changedFiles, line)
	}

	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	idre := regexp.MustCompile(`/([A-F0-9-]+)/content.rtf$`)

	for _, file := range changedFiles {
		matches := idre.FindStringSubmatch(file)
		if matches != nil {
			uuid := matches[1]
			path, ok := paths[uuid]
			if ok {
				fmt.Printf("\n-----\n%s\n-----\n", path)
			} else {
				fmt.Printf("\n-----\nUnknown path: %s\n-----\n", uuid)
			}
			diffargs := append(args, "--word-diff=color", "--no-prefix", file)
			cmd := exec.Command("git", diffargs...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func walkToc(path []string, items []BinderItem) {
	for _, item := range items {
		p := append(path, item.Title)
		paths[item.UUID] = strings.Join(p, " / ")
		walkToc(p, item.Children.BinderItem)
	}
}
