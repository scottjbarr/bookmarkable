package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	errNotFound = errors.New("Not found")

	dbDir         = os.Getenv("HOME") + "/.bookmarkable"
	defaultConfig = dbDir + "/config.json"

	versionFlag = flag.Bool("v", false, "Print version and exit")
	configFile  = flag.String(
		"config",
		defaultConfig,
		"Config file. See config/example.json.dist")
	url  = flag.String("url", "", "URL to bookmark")
	tags = flag.String("tags", "", "\"foo bar\" adds tag \"foo\" and \"bar\"")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Version : %s\n", version)
		fmt.Fprintf(os.Stderr, "Commmit : %s\n", commit)
		fmt.Fprintf(os.Stderr, "Built   : %s\n\n", buildDate)

		flag.PrintDefaults()

		fmt.Fprintln(os.Stderr, "\nUsage:")
		fmt.Fprintln(os.Stderr, "    bookmarkable command -config conf.json")
	}
}

const (
	errorUnparsableConfig = 1
	errorDBCreate         = 2
	errorDBGet            = 4
	errorPageGet          = 8
	errorBookmarkAdd      = 16
	cmdAdd                = "add"
	cmdList               = "list"
	cmdSearch             = "search"
)

func main() {
	// remove the command so that the flags are parsable
	args := os.Args[0:1]
	cmd := os.Args[1]

	fmt.Printf("cmd = %v\n", cmd)

	for _, s := range os.Args[2:] {
		args = append(args, s)
	}

	// command removed, put the args back
	os.Args = args

	flag.Parse()

	if *versionFlag || *url == "" {
		flag.Usage()
		os.Exit(0)
	}

	config, err := parseConfig(*configFile)

	// if cmd == cmdList {
	// 	fmt.Printf("list bookmarks\n")
	// } else if cmd == cmdAdd {
	// 	fmt.Printf("add\n")
	// } else if cmd == cmdSearch {
	// 	fmt.Printf("search\n")
	// }

	fmt.Printf("url = %v\n", *url)

	// os.Exit(0)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(errorUnparsableConfig)
	}

	db := New(config)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(errorDBCreate)
	}

	p, err := NewPage(url)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(errorPageGet)
	}

	fmt.Printf("page = %s\n", p)

	tagArray := strings.Split(*tags, " ")

	b := Bookmark{
		Title:     *p.title,
		URL:       *p.url,
		Tags:      tagArray,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.add(&b); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(errorBookmarkAdd)
	}
}
