FROM golang:1.17-alpine

WORKDIR /app

COPY ./ori $workdir

EXPOSE 9101
CMD ["./ori -c /Users/ian/workdir/cc/goOrigin/config.yaml"]
