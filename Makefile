dev:
	go run main.go

format:
	gofmt -w ./cmd ./internal ./pkg

test:
	go test -v ./...
