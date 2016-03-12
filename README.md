# Bookmarkable

Bookmarkable is a fun little program I've put together to help me
manage my bookmarks across multiple machines, from the command line.

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

## TODO

- Delete bookmark(s)
- Public bookmarks
- Command based docs from command line.
- Error handling

## Licence

The MIT License (MIT)

Copyright (c) 2016 Scott Barr

See [LICENCE.md](LICENCE.md)
