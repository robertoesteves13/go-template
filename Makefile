.PHONY: sql server templ bundle

SASS_FLAGS := --load-path=node_modules/@picocss/pico/scss/ --no-source-map --style compressed
BUN_FLAGS := --minify

server: sql templ bundle
	go build -o web ./cmd/web
sql:
	sqlc generate
templ:
	templ generate
bundle:
	cd cmd/web; bun install
	cd cmd/web; bun build ./index.ts --outfile index.js $(BUN_FLAGS)
	cd cmd/web; sass scss/index.scss:global.css $(SASS_FLAGS)
