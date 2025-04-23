FROM golang:1.24.2 AS builder


RUN apt-get update && apt-get install -y \
    ca-certificates \
    build-essential \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

ENV GOOS=linux
ENV CGO_ENABLED=1

WORKDIR /st
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -a -installsuffix cgo -o ./bin/st ./service-todo/

FROM alpine:latest

RUN apk add --update --no-cache ca-certificates librdkafka libc6-compat

WORKDIR /st

COPY --from=builder /st/bin/st .

EXPOSE 9001
CMD [ "./st" ]