FROM golang:1.14.2 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build/

COPY . .
#ADD go.mod .
#ADD go.sum .
#RUN go mod download

RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /app/app main.go

FROM ubuntu:18.04

RUN apt-get update && apt-get install -y locales ca-certificates tzdata && mkdir -p /app/etc/
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/app /app/app
EXPOSE 8888
CMD ["./app"]