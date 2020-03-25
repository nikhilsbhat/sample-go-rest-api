#first stage - builder
FROM golang:alpine3.11 as builder
RUN go get -u github.com/nikhilsbhat/go-api-sample
COPY . /app
WORKDIR /app
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o sample-api
#second stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/sample-api sample-api
CMD ["./sample-api"]