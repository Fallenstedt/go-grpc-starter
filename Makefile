
VERSION         :=      $(shell cat ./VERSION)
IMAGE_NAME      :=      fallenstedt/test-grpc

all: install proto


install:
	echo "Installing go modules..." && \
	go mod download && \
	echo "Completed" 

proto: 
	echo "Building proto definitions..." && \
	protoc --go_out=plugins=grpc:. greet/greetpb/*.proto && \
	echo "Completed"

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
