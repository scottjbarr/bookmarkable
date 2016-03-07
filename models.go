package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strings"
	"time"
)

func (db *DB) getClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: db.config.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return github.NewClient(tc)
}

func New(config *Config) *DB {
	db := &DB{
		config:    config,
		bookmarks: make([]*Bookmark, 0),
		gistname:  "bookmarkable",
	}

	db.client = db.getClient()
	db.load()

	fmt.Printf("New content = %v\n", *db.gist.Files[db.filename()].Content)
	return db
}

type DB struct {
	config    *Config
	client    *github.Client
	gistname  string
	gist      *github.Gist
	bookmarks []*Bookmark
}

func (db *DB) load() error {
	gists, _, err := db.client.Gists.List("", nil)

	if err != nil {
		return err
	}

	index := -1
	for i, g := range gists {
		if g.Description != nil && *g.Description == db.gistname {
			index = i
		}
	}

	if index > -1 {
		fmt.Printf("loaded existing gist\n")

		id := gists[index].ID
		g, _, err := db.client.Gists.Get(*id)

		if err != nil {
			return err
		}

		db.gist = g

		fmt.Printf("initial content load = %+v\n", *db.gist.Files[db.filename()].Content)

		var ary []*Bookmark

		bytes := []byte(*db.gist.Files[db.filename()].Content)
		err = json.Unmarshal(bytes, &ary)

		if err != nil {
			panic(err)
			return err
		}

		db.bookmarks = ary

		return nil
	}

	return db.create()
}

func (db *DB) create() error {
	var bookmarks = make([]*Bookmark, 0)
	db.bookmarks = bookmarks

	if err := db.updateContent(); err != nil {
		return err
	}

	g, _, err := db.client.Gists.Create(db.gist)

	if err != nil {
		return err
	}

	db.gist = g

	return nil
}

func (db *DB) filename() github.GistFilename {
	var filename github.GistFilename
	filename = "bookmarkable.json"

	return filename
}

func (db *DB) updateContent() error {
	var bytes []byte
	var err error

	if bytes, err = json.MarshalIndent(db.bookmarks, "", "  "); err != nil {
		return err
	}

	content := string(bytes)
	fmt.Printf("content = %v\n", content)

	name := string(db.filename())
	files := make(map[github.GistFilename]github.GistFile, 0)

	files[db.filename()] = github.GistFile{
		Filename: &name,
		Content:  &content,
	}

	db.gist.Files = files // [db.filename()] = file

	fmt.Printf("gist.content = %v\n", *db.gist.Files[db.filename()].Content)

	return nil
}

func (db *DB) add(bookmark *Bookmark) error {
	db.bookmarks = append(db.bookmarks, bookmark)

	db.updateContent()

	var g *github.Gist
	var err error

	g, _, err = db.client.Gists.Edit(*db.gist.ID, db.gist)

	if err != nil {
		return err
	}

	db.gist = g

	return nil
}

type Bookmark struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type page struct {
	url   *string
	title *string
}

func (p *page) String() string {
	return fmt.Sprintf("url:%s : title:%v", *p.url, *p.title)
}

func NewPage(url *string) (*page, error) {
	p := &page{
		url: url,
	}

	doc, err := goquery.NewDocument(*p.url)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(doc.Find("head title").Text())
	p.title = &title

	return p, nil
}
