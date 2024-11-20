FROM golang:1.20-alpine

WORKDIR /app

COPY ./record_exporter $workdir

EXPOSE 9101
CMD ["./record_exporter -c /Users/ian/workdir/cc/goOrigin/config.yaml"]
