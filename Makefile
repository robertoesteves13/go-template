BUN_FLAGS := --minify

GO_FILES := $(shell find . -type f -name "*.go")
TEMPL_FILES := $(shell find . -type f -name "*.templ")
TEMPL_GEN_FILES := $(TEMPL_FILES:.templ=_templ.go)

server: $(GO_FILES) templ cmd/web/services/index.js cmd/web/services/global.css internal/database/db.go
	@echo "SERVER: Building web project"
	@go build -o wserver ./cmd/web

internal/database/db.go: schema.sql query.sql sqlc.yml
	@echo "SQLC: Generating queries"
	@sqlc generate

templ: $(TEMPL_GEN_FILES)

%_templ.go: %.templ
	@echo "TEMPL $<"
	@templ generate -f $<

cmd/web/services/index.js: cmd/web/index.ts
	@echo "BUN: Bundle javascript"
	@cd cmd/web; bun install
	@cd cmd/web; bun build ./index.ts --outfile services/index.js $(BUN_FLAGS)

cmd/web/services/global.css: cmd/web/uno.config.ts
	@echo "UNO: Bundle CSS"
	@cd cmd/web; bun x unocss -c uno.config.ts
	@cd cmd/web; bun x uglifycss global.css --output services/global.css
	@cd cmd/web; rm global.css

clean:
	rm -rf internal/database
	rm -f cmd/web/templates/*_templ.go
	rm -f cmd/web/services/global.css
	rm -f cmd/web/services/index.js
	rm -f wserver

# Default target
.DEFAULT_GOAL := server
