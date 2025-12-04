SHELL := /usr/bin/env bash

PLUGIN_ID := com.mattermost.ai-suite
PLUGIN_VERSION ?= 0.1.0
PLUGIN_PACKAGE := ai-suite

GO ?= go
NPM ?= npm
MMCTL ?= mmctl
BUILD_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
BUILD_DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -X 'main.version=$(PLUGIN_VERSION)' -X 'main.buildHash=$(BUILD_HASH)' -X 'main.buildDate=$(BUILD_DATE)'

DIST_DIR := dist
SERVER_DIST := server/dist
WEBAPP_DIST := webapp/dist
SERVER_TARGETS := \
	$(SERVER_DIST)/plugin-linux-amd64 \
	$(SERVER_DIST)/plugin-linux-arm64 \
	$(SERVER_DIST)/plugin-darwin-amd64 \
	$(SERVER_DIST)/plugin-darwin-arm64 \
	$(SERVER_DIST)/plugin-windows-amd64.exe

.PHONY: help all clean server webapp bundle deploy watch stop-watch lint

help:
	@echo "Mattermost AI Productivity Suite Plugin"
	@echo ""
	@echo "Common targets:"
	@echo "  make server    - Build Go plugin binaries for all major targets"
	@echo "  make webapp    - Install dependencies & build the React bundle"
	@echo "  make bundle    - Produce a distributable plugin tarball"
	@echo "  make deploy    - Deploy the built bundle into a running server (set MM_SERVER_PATH)"
	@echo "  make watch     - Run the webapp in watch mode"
	@echo "  make clean     - Remove build artifacts"

all: bundle

$(SERVER_DIST):
	mkdir -p $@

$(SERVER_DIST)/plugin-linux-amd64: server/*.go | $(SERVER_DIST)
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $@ ./server

$(SERVER_DIST)/plugin-linux-arm64: server/*.go | $(SERVER_DIST)
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $@ ./server

$(SERVER_DIST)/plugin-darwin-amd64: server/*.go | $(SERVER_DIST)
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $@ ./server

$(SERVER_DIST)/plugin-darwin-arm64: server/*.go | $(SERVER_DIST)
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $@ ./server

$(SERVER_DIST)/plugin-windows-amd64.exe: server/*.go | $(SERVER_DIST)
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $@ ./server

server: $(SERVER_TARGETS)

webapp:
	cd webapp && $(NPM) install
	cd webapp && $(NPM) run build

bundle: clean-dist server webapp
	mkdir -p $(DIST_DIR)/$(PLUGIN_ID)
	cp plugin.json $(DIST_DIR)/$(PLUGIN_ID)/
	cp -r $(SERVER_DIST) $(DIST_DIR)/$(PLUGIN_ID)/server
	cp -r $(WEBAPP_DIST) $(DIST_DIR)/$(PLUGIN_ID)/webapp
	cd $(DIST_DIR) && tar -czf $(PLUGIN_PACKAGE)-$(PLUGIN_VERSION).tar.gz $(PLUGIN_ID)
	@echo "Bundle created at $(DIST_DIR)/$(PLUGIN_PACKAGE)-$(PLUGIN_VERSION).tar.gz"

deploy: bundle
ifndef MM_SERVER_PATH
	$(error MM_SERVER_PATH environment variable must point to the root of your Mattermost server installation)
endif
	$(eval BUNDLE_PATH := $(DIST_DIR)/$(PLUGIN_PACKAGE)-$(PLUGIN_VERSION).tar.gz)
	@echo "Deploying $(BUNDLE_PATH) to $(MM_SERVER_PATH)/plugins"
	mkdir -p "$(MM_SERVER_PATH)/plugins"
	tar -xzf "$(BUNDLE_PATH)" -C "$(MM_SERVER_PATH)/plugins"
	@echo "Deployment complete. Enable the plugin via System Console → Plugins."

watch:
	cd webapp && $(NPM) install
	cd webapp && $(NPM) run watch

stop-watch:
	pkill -f "npm run watch" || true

lint:
	cd webapp && $(NPM) run lint

clean-dist:
	rm -rf $(DIST_DIR)
	rm -rf $(SERVER_DIST)
	rm -rf $(WEBAPP_DIST)

clean: clean-dist
	rm -rf webapp/node_modules

