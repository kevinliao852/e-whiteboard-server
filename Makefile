dev:
	go run main.go

test:
	go test -v ./...

e2e-test:
	go test -v ./e2e/...
