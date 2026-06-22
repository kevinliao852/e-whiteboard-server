dev:
	go run ./cmd/...

format:
	gofmt -w ./cmd ./internal ./pkg

test:
	go test -v ./...
