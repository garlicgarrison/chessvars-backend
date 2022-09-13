.PHONY: gen build run test run-docker test-docker

IMAGE_NAME:=garlicgarrison/chessvars-backend

gen:
	go run github.com/99designs/gqlgen generate

build: gen
	rm -rf bin
	CGO_ENABLED=0 go build -v \
		-o bin/ \
		./...

run: gen
	go run ./cmd/backend/...

test:
	go clean -testcache
	go test ./testing/... -v

build-docker:
	docker build \
		-t $(IMAGE_NAME) \
		--progress plain \
		--no-cache \
		.

run-docker:
	docker compose up --detach --build

test-docker: run-docker test

PORT:=8080