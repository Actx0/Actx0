FROM alpine:3.22

RUN apk add --no-cache postgresql16-client

ENTRYPOINT ["psql"]
