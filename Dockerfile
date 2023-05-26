FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:3.18.0
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .

EXPOSE 8080
CMD ["./main"]


