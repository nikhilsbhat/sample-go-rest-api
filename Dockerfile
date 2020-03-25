#first stage - builder
FROM golang:alpine3.11 as builder
COPY . /app
WORKDIR /app
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o go-api-sample
#second stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/go-api-sample go-api-sample
CMD ["./go-api-sample"]