package bookmarkable

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"strings"
	"time"
)

// DB is the front end storage model.
type db struct {
	config         *config
	configFileName *string
	bookmarks      []*Bookmark
	storage        *gistStore
	filename       string
}

// New creates and returns a DB from a config file.
func New(configFileName *string, dbDir string) (*db, error) {
	db := &db{
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

// Sync updates local bookmarks with the remote copy.
//
// This does not merge bookmarks so anything you have managed to store
// locally will be overwritten.
func (db *db) Sync() error {
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

func (db *db) writeConfig() error {
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

func (db *db) writeBookmarks() error {
	d1 := []byte(*db.storage.content)
	return ioutil.WriteFile(db.filename, d1, 0644)
}

// Search returns Bookmarks structs that match the given phrase.
func (db *db) Search(phrase string) []*Bookmark {
	var results []*Bookmark

	ary, err := db.GetBookmarks()

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

// Add adds a new Bookmark to the collection, and updates storage.
func (db *db) Add(url string, tags []string) error {
	p, err := newPage(&url)

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

	bookmarks, err := db.GetBookmarks()
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

// GetBookmarks returns all locally stored Bookmark structs.
func (db *db) GetBookmarks() ([]*Bookmark, error) {
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

// Bookmark represents a single Bookmark
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

func newPage(url *string) (*page, error) {
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
