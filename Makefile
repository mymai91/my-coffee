# .PHONY make sure command always runs, even if a file with the same name exists
.PHONY: all test lint fmt build build-coffee-cli build-brew-svc build-menu-svc build-barista-cli clean gen proto tidy cover start-services stop-services run-coffeecli run-baristacli

all: build

gen: proto

proto:
	bash ./scripts/genproto.sh

tidy:
	go mod tidy

install-tools:
	# Install protoc plugins used by buf generation
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

run-menusvc:
	go run ./cmd/menusvc

run-menusvc-prod:
	ENV=production go run ./cmd/menusvc

run-brewsvc:
	go run ./cmd/brewsvc

run-brewsvc-prod:
	ENV=production go run ./cmd/brewsvc

build-brewsvc:
	go build ./cmd/brewsvc/

run-coffeecli:
	go run ./cmd/coffeecli

docker-db-dev-rm:
	docker compose rm dev-db -s -f -v
docker-db-dev-up:
	docker compose up dev-db -d

# Migrations
migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-status:
	go run cmd/migrate/main.go version

migrate-create:
	@echo "Usage: make migrate-create name=create_users_table"
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(name)