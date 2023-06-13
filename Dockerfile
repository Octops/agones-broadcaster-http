FROM golang:1.20 AS builder

WORKDIR /go/src/github.com/Octops/agones-broadcaster-http

COPY . .

RUN make build && chmod +x /go/src/github.com/Octops/agones-broadcaster-http/bin/broadcaster-http

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /go/src/github.com/Octops/agones-broadcaster-http/bin/broadcaster-http /app/

ENTRYPOINT ["./broadcaster-http"]