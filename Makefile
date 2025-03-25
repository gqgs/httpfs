
build: generate
	go build ./cmd/client
	go build ./cmd/server

generate:
	go generate ./...
