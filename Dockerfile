# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/yfedoruck/webchat

# Build the webchat command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get -u github.com/gorilla/websocket
RUN go get -u golang.org/x/oauth2
RUN go get -u golang.org/x/oauth2/facebook
RUN go get -u golang.org/x/oauth2/google

RUN go install github.com/yfedoruck/webchat

# Run the webchat command by default when the container starts.
ENTRYPOINT /go/bin/webchat

# Document that the service listens on port 8080.
EXPOSE 8080
