FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN GOPROXY=https://goproxy.cn,direct go mod download
COPY . .
RUN go build -o shortlink shortlink.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/shortlink .
COPY --from=builder /app/etc ./etc

EXPOSE 8888
CMD ["./shortlink", "-f", "etc/shortlink-api.yaml"]
