package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"strings"
	"time"
)

type DB struct {
	config    *Config
	bookmarks []*Bookmark
	storage   *gistStore
}

func New(config *Config) *DB {
	db := &DB{
		config:    config,
		bookmarks: make([]*Bookmark, 0),
	}

	db.storage = newGistStore(config.ID, &config.Token)

	if err := db.storage.init(); err != nil {
		panic(err)
	}

	d1 := []byte(*db.storage.content)
	filename := dbDir + "/" + *db.storage.file.Filename

	if err := ioutil.WriteFile(filename, d1, 0644); err != nil {
		panic(err)
	}

	return db
}

// func (db *DB) create() error {
// 	var bookmarks = make([]*Bookmark, 0)
// 	db.bookmarks = bookmarks

// 	if err := db.updateContent(); err != nil {
// 		return err
// 	}

// 	g, _, err := db.client.Gists.Create(db.gist)

// 	if err != nil {
// 		return err
// 	}

// 	db.gist = g

// 	return nil
// }

// func (db *DB) filename() github.GistFilename {
// 	var filename github.GistFilename
// 	filename = "bookmarkable.json"

// 	return filename
// }

// func (db *DB) updateContent() error {
// 	var bytes []byte
// 	var err error

// 	if bytes, err = json.MarshalIndent(db.bookmarks, "", "  "); err != nil {
// 		return err
// 	}

// 	content := string(bytes)
// 	fmt.Printf("content = %v\n", content)

// 	name := string(db.filename())
// 	files := make(map[github.GistFilename]github.GistFile, 0)

// 	files[db.filename()] = github.GistFile{
// 		Filename: &name,
// 		Content:  &content,
// 	}

// 	db.gist.Files = files // [db.filename()] = file

// 	fmt.Printf("gist.content = %v\n", *db.gist.Files[db.filename()].Content)

// 	return nil
// }

func (db *DB) add(bookmark *Bookmark) error {
	db.bookmarks = append(db.bookmarks, bookmark)

	var bytes []byte
	var err error

	if bytes, err = json.MarshalIndent(db.bookmarks, "", "  "); err != nil {
		return err
	}

	content := string(bytes)
	return db.storage.update(&content)
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
