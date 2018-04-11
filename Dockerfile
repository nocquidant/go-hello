# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# create the user
RUN useradd -r -s /bin/false helloworld
# Go get and build
RUN go get github.com/nocquidant/go-hello
RUN go install github.com/nocquidant/go-hello

# Run the service
ENTRYPOINT /go/bin/go-hello

# Document that the service listens on port 8080.
EXPOSE 8484
