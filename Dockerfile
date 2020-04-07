# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine

WORKDIR /go/src/webchat

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/webchat github.com/yfedoruck/webchat

# Run the webchat command by default when the container starts.
CMD ["/bin/webchat"]

# Document that the service listens on port 8080.
EXPOSE 8080
