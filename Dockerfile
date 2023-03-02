FROM golang:1.20 AS builder
COPY . /workspace
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /output/go-http-echo /workspace/main.go

FROM scratch
COPY --from=builder /output/go-http-echo /go-http-echo
ENTRYPOINT ["/go-http-echo"]