.PHONY: sql server templ bundle

BUN_FLAGS := --minify

server: sql templ bundle
	@echo "SERVER: Building web project"
	@go build -o wserver ./cmd/web

sql:
	@echo "SQLC: Generating queries"
	@sqlc generate

templ:
	@echo "TEMPL: Generating templates"
	@templ generate

bundle:
	@echo "TS: Bundle javascript"
	@cd cmd/web; bun install
	@cd cmd/web; bun build ./index.ts --outfile services/index.js $(BUN_FLAGS)
	@echo "UNO: Bundle CSS"
	@cd cmd/web; bunx unocss -c uno.config.ts
	@cd cmd/web; bunx uglifycss global.css --output services/global.css

clean:
	rm -f cmd/web/services/global.css
	rm -f cmd/web/services/index.js
	rm -f wserver
