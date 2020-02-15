default: test

build:
	go build -v ./...

lint:
	go vet ./...

test:
	go test ./... -v
