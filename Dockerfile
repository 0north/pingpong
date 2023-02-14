ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-alpine AS builder

RUN mkdir -p /api
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./app ./src/main.go

FROM alpine:latest

COPY --from=builder /api/app .

EXPOSE 8080

ENTRYPOINT ["./app"]
