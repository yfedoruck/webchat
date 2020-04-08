# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine  AS builder

WORKDIR /go/src/webchat

RUN apk add --no-cache curl \
    bash \
    git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/webchat github.com/yfedoruck/webchat/cmd/chat

FROM alpine:3.11
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/webchat /go/src/webchat
COPY --from=builder /bin/webchat /bin/webchat

# Run the webchat command by default when the container starts.
CMD ["/bin/webchat"]

# Document that the service listens on port 8080.
EXPOSE 8080
