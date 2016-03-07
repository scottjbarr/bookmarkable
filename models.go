package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/github"
	"golang.org/x/net/html"
	"golang.org/x/oauth2"
	"net/http"
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
	// db.load()

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
		db.gist = &gists[index]
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

func (db *DB) updateContent() error {
	var bytes []byte
	var err error

	if bytes, err = json.MarshalIndent(db.bookmarks, "", "  "); err != nil {
		return err
	}

	content := string(bytes)
	fmt.Printf("content = %v\n", content)

	var filename github.GistFilename
	filename = "bookmarkable.json"

	if db.gist == nil {
		fmt.Printf("updateContent (nil gist)\n")
		name := string(filename)
		files := make(map[github.GistFilename]github.GistFile, 0)

		files[filename] = github.GistFile{
			Filename: &name,
			Content:  &content,
		}

		public := false
		db.gist = &github.Gist{
			Description: &db.gistname,
			Files:       files,
			Public:      &public,
		}
	} else {
		fmt.Printf("updateContent (existing gist)\n")
		m := db.gist.Files[filename]
		m.Content = &content
	}

	fmt.Printf("gist.content = %v\n", *db.gist.Files[filename].Content)

	return nil
}

func (db *DB) add(bookmark *Bookmark) error {
	db.bookmarks = append(db.bookmarks, bookmark)

	// fmt.Printf("id = %v\n", *db.gist.ID)
	// fmt.Printf("gist = %v\n", *db.gist)

	db.updateContent()

	var g *github.Gist
	var err error

	if db.gist.ID != nil {
		g, _, err = db.client.Gists.Edit(*db.gist.ID, db.gist)
	} else {
		g, _, err = db.client.Gists.Create(db.gist)
	}

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

	title := doc.Find("head title").Text()
	p.title = &title

	return p, nil
}

func (p *page) load() error {
	response, err := http.Get(*p.url)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	// body, err := ioutil.ReadAll(response.Body)

	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("body = %v\n", body)

	// return &body, nil

	// return response.Body, nil

	// b := response.Body

	z := html.NewTokenizer(response.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return nil
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			if t.Data != "title" {
				continue
			}

			// fmt.Printf("title z = %+v\n", z)
			// fmt.Printf("title t = %+v\n", t)
			// fmt.Printf("raw = %+v\n", string(z.Raw()))
			// fmt.Printf("next = %+v\n", tt.FirstChild())
			// Extract the href value, if there is one
			// ok, url := getHref(t)
			getTitle(&t)
			// if !ok {
			// 	continue
			// }

			// // Make sure the url begines in http**
			// hasProto := strings.Index(url, "http") == 0
			// if hasProto {
			// 	ch <- url
			// }
		}
	}

	return nil
}

func getTitle(t *html.Token) string {
	fmt.Printf("t = %+v\n", *t)
	fmt.Printf("t = %+v\n", *t)
	return ""
}

// func (p *page) parse(bytes *[]byte) (string, error) {
// 	r := io.NewReader(*bytes)
// 	d := html.NewTokenizer(r)
// 	for {
// 		// token type
// 		tokenType := d.Next()
// 		if tokenType == html.ErrorToken {
// 			return
// 		}
// 		token := d.Token()
// 		switch tokenType {
// 		case html.StartTagToken: // <tag>
// 			// type Token struct {
// 			//     Type     TokenType
// 			//     DataAtom atom.Atom
// 			//     Data     string
// 			//     Attr     []Attribute
// 			// }
// 			//
// 			// type Attribute struct {
// 			//     Namespace, Key, Val string
// 			// }
// 		case html.TextToken: // text between start and end tag
// 		case html.EndTagToken: // </tag>
// 		case html.SelfClosingTagToken: // <tag/>

// 		}
// 	}
// }
