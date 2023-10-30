FROM golang:1.16-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o my-exporter .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/my-exporter .

EXPOSE 9115

CMD ["./my-exporter"]