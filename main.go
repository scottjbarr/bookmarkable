package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	dbDir         = os.Getenv("HOME") + "/.bookmarkable"
	defaultConfig = dbDir + "/config.json"

	versionFlag = flag.Bool("v", false, "Print version and exit")
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
	cmdSync   = "sync"
	cmdAdd    = "add"
	cmdList   = "list"
	cmdSearch = "search"
)

func main() {
	// remove the command so that the flags are parsable
	args := os.Args[0:1]
	cmd := os.Args[1]

	for _, s := range os.Args[2:] {
		args = append(args, s)
	}

	// command removed, put the args back
	os.Args = args

	flag.Parse()

	if *versionFlag {
		flag.Usage()
		os.Exit(0)
	}

	configFile := &defaultConfig
	db, err := New(configFile)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	if cmd == cmdSync {
		if err := db.sync(); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(2)
		}
	} else if cmd == cmdSearch {
		results := db.search(os.Args[1])
		printBookmarks(results)
	} else if cmd == cmdList {
		results, _ := db.getBookmarks()
		printBookmarks(results)
	} else if cmd == cmdList {
		url := os.Args[1]
		tags := os.Args[2:]
		if err := db.add(url, tags); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(4)
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
