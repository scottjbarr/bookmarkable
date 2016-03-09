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
	filename  string
}

func New(config *Config) *DB {
	db := &DB{
		config:    config,
		bookmarks: make([]*Bookmark, 0),
	}
	db.filename = dbDir + "/bookmarkable.json"
	db.storage = newGistStore(config.ID, &config.Token)

	return db
}

func (db *DB) sync() {
	if err := db.storage.init(); err != nil {
		panic(err)
	}

	d1 := []byte(*db.storage.content)

	if err := ioutil.WriteFile(db.filename, d1, 0644); err != nil {
		panic(err)
	}
}

func (db *DB) search(phrase string) []*Bookmark {
	results := make([]*Bookmark, 0)

	ary, err := db.getBookmarks()

	if err != nil {
		panic(err)
	}

	for _, b := range ary {
		if b.matches(phrase) {
			results = append(results, b)
		}
	}

	return results
}

func (db *DB) add(url string, tags []string) error {
	p, err := NewPage(&url)

	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	fmt.Printf("page = %s\n", p)

	b := Bookmark{
		Title:     *p.title,
		URL:       *p.url,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	fmt.Printf("b = %+v\n", b)

	db.bookmarks = append(db.bookmarks, &b)

	var bytes []byte

	if bytes, err = json.MarshalIndent(db.bookmarks, "", "  "); err != nil {
		return err
	}

	content := string(bytes)
	fmt.Printf("content = %v\n", content)
	// return db.storage.update(&content)
	return nil
}

func (db *DB) getBookmarks() ([]*Bookmark, error) {
	if len(db.bookmarks) == 0 {
		bytes, err := ioutil.ReadFile(db.filename)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(bytes, &db.bookmarks); err != nil {
			return nil, err
		}
	}

	return db.bookmarks, nil
}

type Bookmark struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Bookmark) matches(phrase string) bool {
	phrase = strings.ToLower(phrase)

	if strings.Contains(strings.ToLower(b.Title), phrase) {
		return true
	}

	for _, tag := range b.Tags {
		if strings.Contains(tag, phrase) {
			return true
		}
	}

	return false
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
