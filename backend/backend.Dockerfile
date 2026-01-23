FROM golang:latest

WORKDIR /app

RUN go install \
    golang.org/x/tools/gopls@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

ENV GO111MODULE=on \ 
    GOPATH=/go \
    PATH=$PATH:/go/bin