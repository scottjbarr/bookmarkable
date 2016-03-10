package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
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

	// fmt.Printf("cmd = %v\n", cmd)

	for _, s := range os.Args[2:] {
		args = append(args, s)
	}

	// command removed, put the args back
	os.Args = args

	flag.Parse()

	if *versionFlag { // }|| *url == "" {
		flag.Usage()
		os.Exit(0)
	}

	config, err := parseConfig(*configFile)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(errorUnparsableConfig)
	}

	db := New(config)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(errorDBCreate)
	}

	if cmd == "sync" {
		db.sync()
	} else if cmd == "search" {
		results := db.search(os.Args[1])
		printBookmarks(results)
	} else if cmd == "list" {
		results, _ := db.getBookmarks()
		printBookmarks(results)
	} else if cmd == "add" {
		// fmt.Printf("adding\n")
		url := os.Args[1]
		tags := os.Args[2:]
		// fmt.Printf("url = %v tags = %v\n", url, tags)
		if err := db.add(url, tags); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(64)
		}
	}
}

func printBookmarks(bookmarks []*Bookmark) {
	for _, b := range bookmarks {
		fmt.Printf("%v\n  %v\n  %v\n  %v\n\n",
			b.Title,
			b.URL,
			b.Tags,
			b.CreatedAt)
	}
}
