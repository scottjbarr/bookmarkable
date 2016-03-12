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
	config         *Config
	configFileName *string
	bookmarks      []*Bookmark
	storage        *gistStore
	filename       string
}

func New(configFileName *string) (*DB, error) {
	db := &DB{
		configFileName: configFileName,
		bookmarks:      make([]*Bookmark, 0),
	}

	c, err := parseConfig(*db.configFileName)

	if err != nil {
		return nil, err
	}

	db.config = c
	db.filename = dbDir + "/bookmarkable.json"
	db.storage = newGistStore(db.config.ID, &db.config.Token)

	return db, nil
}

func (db *DB) sync() error {
	var err error
	if err = db.storage.init(); err != nil {
		return err
	}

	// write the gist id to the config file if we don't already have it
	if db.config.ID == nil {
		db.config.ID = db.storage.gist.ID
		if err = db.writeConfig(); err != nil {
			return err
		}
	}

	return db.writeBookmarks()
}

func (db *DB) writeConfig() error {
	var b []byte
	var err error

	if b, err = json.MarshalIndent(db.config, "", "  "); err != nil {
		return err
	}

	if err = ioutil.WriteFile(*db.configFileName, b, 0644); err != nil {
		return err
	}

	return nil
}

func (db *DB) writeBookmarks() error {
	d1 := []byte(*db.storage.content)
	return ioutil.WriteFile(db.filename, d1, 0644)
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
		return err
	}

	b := Bookmark{
		Title:     *p.title,
		URL:       *p.url,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	bookmarks, err := db.getBookmarks()
	if err != nil {
		return err
	}

	db.bookmarks = append(bookmarks, &b)

	var bytes []byte

	if bytes, err = json.MarshalIndent(db.bookmarks, "", "  "); err != nil {
		return err
	}

	content := string(bytes)
	db.storage.content = &content

	if err := db.storage.update(); err != nil {
		return err
	}

	return db.writeBookmarks()
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
