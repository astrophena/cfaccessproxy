PREFIX  ?= $(HOME)

BIN     = cfaccessproxy
BINDIR  = $(PREFIX)/bin

LDFLAGS = "-s -w"

.PHONY: build install clean help

build: ## Build
	@ CGO_ENABLED=0 go build -o $(BIN) -trimpath -ldflags=$(LDFLAGS)

install: build ## Install
	@ mkdir -m755 -p $(BINDIR) && \
		install -m755 $(BIN) $(BINDIR)

clean: ## Clean
	@ rm -f $(BIN)

help: ## Show help
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[0;32m%-30s\033[0m %s\n", $$1, $$2}'
