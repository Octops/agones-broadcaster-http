FROM golang:1.14 AS builder

WORKDIR /go/src/github.com/Octops/agones-broadcaster-http

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/broadcaster-http .

FROM alpine

WORKDIR /app

COPY --from=builder /go/src/github.com/Octops/agones-broadcaster-http/bin/broadcaster-http /app/

RUN chmod +x broadcaster-http

ENTRYPOINT ["./broadcaster-http"]