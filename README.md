# Bookmarkable

Bookmarkable is a little program I've built to help me manage my
bookmarks across multiple machines, from the command line. Yes, I
could have (and have) used some third party thing, but in the end I
just want simplicity and to be in control of my own data.

I'm using a secret Gist as storage, so you will neeed to setup a Personal Access Token on Github.

## Install

Install the command line program.

```
go get -u github.com/scottjbarr/bookmarkable/cmd/bookmarkable
```

## Setup

### Create a Personal Access Token on Github.

Create a Personal Access Token on Github. The only permission required
is `Create Gist`. Take note of this token after creation as you cannot
view it again. (Worst case, you can regenerate the token.)

### Create a configuration file

Create a config file at `~/.bookmarkable/config.json`.

If your token was `foo`, the config file will need to contain...

```
{
  "token": "foo"
}
```

## Run

Sync your bookmarks.

```
bookmarkable sync
```

Add a bookmark. Tags are optional.

```
bookmarkable add http://example.com tag0 tag1 tag2
```

List bookmarks.

```
bookmarkable list
```

Search bookmarks for `foo`.

```
bookmarkable search foo
```

## TODO?

- Command based docs from command line
- Error handling
- Delete bookmark(s)
- Public bookmarks
- Publish a HTML friendly list of bookmarks
- Alternative storage

## Licence

The MIT License (MIT)

Copyright (c) 2016 Scott Barr

See [LICENCE.md](LICENCE.md)
