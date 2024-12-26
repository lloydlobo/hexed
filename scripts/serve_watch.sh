#!/usr/bin/env bash

# $ chmod +x scripts/serve_watch.sh
# $ ./scripts/serve_watch.sh
#
# Watch for changes in .go, .html, .js (limited to ./js), and .css (limited to ./css)
# Run build and serve commands on change

find . -type f \( \
    -name "*.go" -o \
    -name "*.html" -o \
    -path "./js/*.js" -o \
    -path "./css/*.css" \
\) | entr -crs '
    date --utc
    go run cmd/make/make.go build
    go run cmd/make/make.go serve
'
