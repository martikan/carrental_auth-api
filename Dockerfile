# Stage 1 - BUILD
FROM golang:1.19-alpine3.16 AS builder

WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Stage 2 - RUN
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder app/start.sh .
COPY --from=builder /app/migrate ./migrate
COPY db/migrations ./migrations

EXPOSE 8000
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]