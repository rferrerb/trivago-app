FROM scratch
COPY main /go/bin/main
ENTRYPOINT ["/go/bin/main"]