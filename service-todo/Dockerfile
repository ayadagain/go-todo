FROM public.ecr.aws/docker/library/golang:1.24 as builder

ENV CGO_ENABLED 1
ENV GOFLAGS -mod=vendor
ENV GOOS=linux
ENV GO111MODULE=on


WORKDIR /st

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -a -installsuffix cgo -o ./bin/st ./service-payment/

FROM alpine:latest
WORKDIR /st

COPY --from=builder /st/bin/st .

EXPOSE 9001
CMD [ "./st" ]