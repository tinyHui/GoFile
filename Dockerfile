FROM golang

RUN mkdir -p $GOPATH/src
RUN mkdir -p $GOPATH/pkg
RUN mkdir -p $GOPATH/bin

WORKDIR /tmp
RUN curl https://glide.sh/get | sh

# Add code
ADD cmd /go/src/github.com/tinyhui/GoFile/cmd
ADD fileop /go/src/github.com/tinyhui/GoFile/fileop
ADD router /go/src/github.com/tinyhui/GoFile/router
ADD utils /go/src/github.com/tinyhui/GoFile/utils
ADD glide.yaml /go/src/github.com/tinyhui/GoFile

# Add config file
RUN mkdir /config
ADD config/prod/parameters.yaml /config/

# Storage folder
RUN mkdir /storage

# Build
WORKDIR /go/src/github.com/tinyhui/GoFile
RUN glide up
RUN go build -o $GOPATH/bin/main ./cmd

# Run
EXPOSE 8080
ENV GIN_MODE release
ENV config /config/parameters.yaml
ENTRYPOINT ["/go/bin/main"]
