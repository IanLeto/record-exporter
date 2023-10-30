FROM golang:1.16-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o record-exporter .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/record-exporter .

EXPOSE 9101

CMD ["./record-exporter"]