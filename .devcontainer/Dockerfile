FROM golang:1.19 as build

RUN (which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1
RUN go install -v golang.org/x/tools/gopls@latest
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest
RUN go install -v honnef.co/go/tools/cmd/staticcheck@latest
RUN apt-get update && apt-get install -y git

#ENV http_proxy socks5://192.168.88.1:1080
#ENV https_proxy socks5://192.168.88.1:1080

