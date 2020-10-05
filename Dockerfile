FROM golang:1.14.2-alpine3.11

# Required for fetching the dependencies
RUN apk update && apk add --no-cache git ca-certificates openssh-client curl

ENV GO111MODULE on

WORKDIR $GOPATH/src/github.com/abhinavk1/remote-config-server

COPY go.mod .
COPY go.sum .
COPY ./cmd ./cmd
COPY ./pkg ./pkg

# Download dependencies
RUN go mod download -x

# Build the binaries
WORKDIR $GOPATH/src/github.com/abhinavk1/remote-config-server/cmd/remote-config-server
RUN go build -tags=jsoniter -o /remote-config-server