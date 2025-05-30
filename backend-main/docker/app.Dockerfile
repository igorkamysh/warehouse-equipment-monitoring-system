FROM golang:alpine AS builder 

WORKDIR /build

ADD go.mod .
ADD go.sum .

COPY . .

RUN go build -o app cmd/app/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/app /build/app

COPY config /build/config

CMD ["./app", "--config", "config/local.yml"]
