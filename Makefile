
.PHONY: test
test:
	@go test -race ./...

.PHONY: cover
cover:
	@go test -cover -race ./...

.PHONY: build
build:
	CGO_ENABLED=0 go build ./...
