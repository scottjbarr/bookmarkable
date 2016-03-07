package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	versionFlag = flag.Bool("v", false, "Print version and exit")
	configFile  = flag.String(
		"config",
		"/Users/scott/.bookmarkable/config.json",
		"Config file. See config/example.json.dist")
	url  = flag.String("url", "", "URL to bookmark")
	tags = flag.String("tags", "", "e.g. \"foo bar\" adds tag \"foo\" and \"bar\"")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Version : %s\n", version)
		fmt.Fprintf(os.Stderr, "Commmit : %s\n", commit)
		fmt.Fprintf(os.Stderr, "Built   : %s\n\n", buildDate)

		flag.PrintDefaults()

		fmt.Fprintln(os.Stderr, "\nUsage:")
		fmt.Fprintln(os.Stderr, "    bookmarkable -config conf.json")
	}
}

const (
	errorUnparsableConfig = 1
	errorDBCreate         = 2
	errorDBGet            = 4
	errorPageGet          = 8
	errorBookmarkAdd      = 16
)

func main() {
	flag.Parse()

	if *versionFlag || *url == "" {
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

	// fmt.Printf("gist = %+v\n", *db.gist)

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
