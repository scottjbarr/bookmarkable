# Bookmarkable

Bookmarkable is a fun little program I've put together to help me manage my
bookmarks across multiple machines, from the command line.

I'm using a secret Gist as storage.

## Setup

### Create a Personal Access Token on Github.

Create a Personal Access Token on Github. The only permission required is `Create Gist`. Take note of this token after creation as You cannot view it again.

### Create a configuration file

Create a config file at `~/.bookmarkable/config.json`.

If your token was `foo`, the config file will need to contain...

```
{
  "token": "foo"
}
```

## Install

Install the command line program.

```
go get -u github.com/scottjbarr/bookmarkable/cmd/bookmarkable
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

## TODO

- Public bookmarks
- Command based docs from command line.

## Licence

The MIT License (MIT)

Copyright (c) 2016 Scott Barr

See [LICENCE.md](LICENCE.md)
