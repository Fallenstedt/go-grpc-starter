FROM golang:1.14.6-alpine AS builder
RUN apk update && \ 
    apk add --no-cache wget && \ 
    apk add make && \
    apk add protoc && \
    apk add git

RUN go get -u github.com/golang/protobuf/proto
RUN go get -u github.com/golang/protobuf/protoc-gen-go

# Deps
WORKDIR /app
COPY Makefile .
COPY go.mod .
COPY go.sum .
COPY VERSION . 
RUN make install

# Build
COPY . .
RUN make proto
RUN make build

FROM alpine
COPY --from=builder /app/server /app/

EXPOSE 8080

# # Set the binary as the entrypoint of the container
ENTRYPOINT ["/app/server"]
