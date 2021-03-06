
VERSION         :=      $(shell cat ./VERSION)
IMAGE_NAME      :=      fallenstedt/test-grpc

all: install gen


install:
	echo "Installing go modules..." && \
	go mod download && \
	echo "Completed" 


proto:
	protoc --go_out=gen --go_opt=paths=source_relative \
    --go-grpc_out==plugins=grpc:gen --go-grpc_opt=paths=source_relative \
    greet/proto/*.proto

build:
	go build greet/greet_server/server.go

test:
	go test ./greet/... -v

image:
	docker build -t fallenstedt/test-grpc .

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)

.PHONY: install test fmt build release
