FROM alpine:3.8

WORKDIR /usr/src/app
COPY ./bin/cron ./

CMD ["./cron"]