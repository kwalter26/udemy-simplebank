FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .
RUN apk add --no-cache curl

ARG TARGETPLATFORM
RUN if [ "$TARGETPLATFORM" = "linux/amd64" ]; then ARCHITECTURE=amd64; elif [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then ARCHITECTURE=arm64; else ARCHITECTURE=amd64; fi \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-${ARCHITECTURE}.tar.gz  | tar xvz

FROM alpine:3.18.0

ARG USERNAME=banker
ARG GROUP=banker

# Create the user
RUN addgroup -S $GROUP && adduser -S $USERNAME -G $GROUP

USER $USERNAME

WORKDIR /app
COPY --chmod=544 --chown=$USERNAME:$GROUP --from=builder /app/main ./
COPY --chmod=544 --chown=$USERNAME:$GROUP --from=builder /app/migrate ./
# COPY as a user
COPY --chmod=544 --chown=$USERNAME:$GROUP start.sh ./
COPY --chmod=544 --chown=$USERNAME:$GROUP db/migration/ ./db/migration/

ENV GIN_MODE=release

EXPOSE 8080
CMD ["/app/main"]
#ENTRYPOINT ["/app/start.sh"]


