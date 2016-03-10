# See http://peter.bourgon.org/go-in-production/
GO ?= go

BUILD_DIR = build
DIST_DIR = dist

PROG = bookmarkable
PROG_BUILD = $(BUILD_DIR)/$(PROG)
PROG_DIST = $(DIST_DIR)/$(PROG)

BUILD = $(GO) build -ldflags $(FLAGS)
INSTALL= $(GO) install -ldflags $(FLAGS)

VERSION := `git rev-parse HEAD`
COMMIT := `git rev-parse HEAD`
BUILD_DATE := `date +%Y-%m-%d\ %H:%M`
FLAGS := "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X \"main.buildDate=$(BUILD_DATE)\""

# the default config file for the "run" target
CONFIG_FILE = config/development.json

GO_FILES = `ls *.go | grep -v test | xargs echo`

all: clean build

build:
	$(BUILD) -o $(PROG_BUILD)

install:
	$(INSTALL)

run:
	$(GO) run $(GO_FILES) -config $(CONFIG_FILE) -url http://example.com -tags "foo,bar"

test:
	$(GO) test

clean:
	rm -rf $(BUILD_DIR)

distclean:
	rm -rf $(PROJ_DIST)
