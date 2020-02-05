default: build

build:
	go build github.com/dkinzer/go-linode-ip/cmd/linode-ip

test:
	go test -v github.com/dkinzer/go-linode-ip/cmd/linode-ip
