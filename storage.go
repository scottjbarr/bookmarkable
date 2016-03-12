package bookmarkable

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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

	return err
}

func (gs *gistStore) byID() error {
	gist, _, err := gs.client.Gists.Get(*gs.id)

	if err != nil {
		return err
	}

	gs.gist = gist

	file := gist.Files[*gs.filename]
	gs.file = &file
	gs.content = gs.file.Content

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
}

func (gs *gistStore) update() error {
	files := make(map[github.GistFilename]github.GistFile, 0)

	f := "bookmarkable.json"
	files[*gs.filename] = github.GistFile{
		Filename: &f,
		Content:  gs.content,
	}

	gs.gist = &github.Gist{}
	gs.gist.Files = files

	gist, _, err := gs.client.Gists.Edit(*gs.id, gs.gist)

	if err != nil {
		return err
	}

	gs.gist = gist
	return nil
}
