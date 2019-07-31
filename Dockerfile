# Builder

FROM golang:1.12-alpine as builder

# Setup builder
RUN apk update && apk add git make protobuf protobuf-dev

# Install protoc and related
RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger


# Build app
ENV GO111MODULE on

WORKDIR $GOPATH/src/github.com/migotom/cell-centre-services

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go install ./...

