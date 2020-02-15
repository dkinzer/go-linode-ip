default: test

build:
	go build -v ./...

lint:
	go vet ./... -v

test:
	go test ./... -v
