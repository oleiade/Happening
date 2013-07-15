HAPPENING_PACKAGE := github.com/oleiade/Happening

BUILD_DIR := $(CURDIR)/.gopath

GOPATH ?= $(BUILD_DIR)
export GOPATH

GO_OPTIONS ?=
ifeq ($(VERBOSE), 1)
GO_OPTIONS += -v
endif

GIT_COMMIT = $(shell git rev-parse --short HEAD)
GIT_STATUS = $(shell test -n "`git status --porcelain`" && echo "+CHANGES")

NO_MEMORY_LIMIT ?= 0
export NO_MEMORY_LIMIT

BUILD_OPTIONS = -ldflags "-X main.GIT_COMMIT $(GIT_COMMIT)$(GIT_STATUS) -X main.NO_MEMORY_LIMIT $(NO_MEMORY_LIMIT)"

SRC_DIR := $(GOPATH)/src

HAPPENING_DIR := $(SRC_DIR)/$(HAPPENING_PACKAGE)
HAPPENING_MAIN := $(HAPPENING_DIR)/happening

HAPPENING_BIN_RELATIVE := bin/happening
HAPPENING_BIN := $(CURDIR)/$(HAPPENING_BIN_RELATIVE)

.PHONY: all clean test

all: $(HAPPENING_BIN)

$(HAPPENING_BIN): $(HAPPENING_DIR)
	# Proceed to happening build
	@(mkdir -p  $(dir $@))
	@(cd $(HAPPENING_MAIN); go get $(GO_OPTIONS); go build $(GO_OPTIONS) $(BUILD_OPTIONS) -o $@)
	@echo $(HAPPENING_BIN_RELATIVE) is created.

$(HAPPENING_DIR):
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR)/ $@

clean:
ifeq ($(GOPATH), $(BUILD_DIR))
	@rm -rf $(BUILD_DIR)
else ifneq ($(HAPPENING_DIR), $(realpath $(HAPPENING_DIR)))
	@rm -f $(HAPPENING_DIR)
endif

test: all
	@(go get "github.com/stretchr/testify/assert")
	@(cd $(HAPPENING_DIR); sudo -E go test $(GO_OPTIONS))

fmt:
	@gofmt -s -l -w .
