### Description: Dockerfile for sample-go-rest-api
FROM alpine:3.14

COPY sample-go-rest-api /

# Starting
CMD [ "/sample-go-rest-api" ]