package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// type storage interface {
// 	init()
// 	sync()
// 	save()
// }

type gistStore struct {
	id       *string
	token    *string
	name     *string
	client   *github.Client
	content  *string
	filename *github.GistFilename
	gist     *github.Gist
	file     *github.GistFile
}

func (gs *gistStore) getClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *gs.token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return github.NewClient(tc)
}

func newGistStore(id *string, token *string) *gistStore {
	name := "bookmarkable"
	gs := gistStore{
		id:    id,
		name:  &name,
		token: token,
	}

	var filename github.GistFilename
	filename = "bookmarkable.json"
	gs.filename = &filename

	gs.client = gs.getClient()

	return &gs
}

func (gs *gistStore) init() error {
	var err error

	err = gs.retrieve()

	if err != nil {
		return err
	}

	if gs.gist != nil {
		return nil
	}

	// gs.create()

	return err
}

func (gs *gistStore) create() error {
	return nil
}

func (gs *gistStore) byID() error {
	fmt.Printf("byID()\n")
	gist, _, err := gs.client.Gists.Get(*gs.id)

	if err != nil {
		return err
	}

	gs.gist = gist

	file := gist.Files[*gs.filename]
	gs.file = &file
	// gs.content = gist.Files[*gs.filename].Content
	gs.content = gs.file.Content
	fmt.Printf("  gist content = %v\n", *gs.content)

	return nil
}

func (gs *gistStore) retrieve() error {
	if gs.id != nil {
		return gs.byID()
	}

	gists, _, err := gs.client.Gists.List("", nil)

	if err != nil {
		return err
	}

	for _, g := range gists {
		if g.Description != nil && *g.Description == *gs.name {
			gs.id = g.ID
		}
	}

	// not found
	if gs.id == nil {
		return fmt.Errorf("Not found")
	}

	return gs.byID()

	// if err != nil {
	// 	return err
	// }

	// db.gist = g

	// fmt.Printf("initial content load = %+v\n", *g.Files[db.filename()].Content)

	// var ary []*Bookmark

	// bytes := []byte(*db.gist.Files[db.filename()].Content)
	// err = json.Unmarshal(bytes, &ary)

	// if err != nil {
	// 	panic(err)
	// 	return err
	// }

	// db.bookmarks = ary
	// file := gs.gist.Files[*gs.filename]
	// return nil // file.Content, nil
}

func (gs *gistStore) update(content *string) error {
	// g, _, err = gs.client.Gists.Edit(gs.ID, gs.gist)
	return nil
}
