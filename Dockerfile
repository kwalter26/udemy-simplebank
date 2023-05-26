FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz  | tar xvz

FROM alpine:3.18.0

ARG USERNAME=banker
ARG GROUP=banker

# Create the user
RUN addgroup -S $GROUP && adduser -S $USERNAME -G $GROUP

USER $USERNAME

WORKDIR /app
COPY --chown=$USERNAME:$GROUP --from=builder /app/main .
COPY --chown=$USERNAME:$GROUP --from=builder /app/migrate ./migrate
# COPY as a user
COPY --chown=$USERNAME:$GROUP start.sh .
COPY --chown=$USERNAME:$GROUP app.env .
COPY --chown=$USERNAME:$GROUP db/migration ./migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]


