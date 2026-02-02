# .PHONY make sure command always runs, even if a file with the same name exists
.PHONY: all test lint fmt build build-coffee-cli build-brew-svc build-menu-svc build-barista-cli clean gen proto tidy cover start-services stop-services run-coffeecli run-baristacli

all: build

gen: proto

proto:
	bash ./scripts/genproto.sh

install-tools:
	# Install protoc plugins used by buf generation
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
