FROM alpine:3.8

# Root Certificates needed for making https/ssl requests
# Bash needed for Kubernetes Dashboard Shell access
RUN apk update && \
  apk add ca-certificates && \
  update-ca-certificates

WORKDIR /usr/src/app
COPY ./bin/cron ./

CMD ["./cron"]