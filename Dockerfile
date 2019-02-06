# SETP 1 build executable binary
FROM golang:alpine AS builder
# Install git to download dependencies
RUN apk update && apk add --no-cache git
COPY main.go .

RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/main



#SETP 2 build a
FROM scratch
COPY --from=builder /go/bin/main /go/bin/main
ENTRYPOINT ["/go/bin/main"]