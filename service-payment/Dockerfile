FROM public.ecr.aws/docker/library/golang:1.24 AS builder

ENV CGO_ENABLED 1
ENV GOFLAGS -mod=vendor
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

WORKDIR /sp

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -a -installsuffix cgo -o ./bin/sp ./service-payment/

FROM alpine:latest
WORKDIR /sp

COPY --from=builder /sp/bin/sp .

EXPOSE 9001
CMD [ "./sp" ]