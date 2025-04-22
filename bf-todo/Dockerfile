FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /bff

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o ./bin/bff ./bf-todo/

FROM alpine:latest
WORKDIR /bff

COPY --from=builder /bff/bin/bff .

EXPOSE 8080
CMD [ "./bff" ]