FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY . .

RUN go build -v -ldflags="-w -s" -o /app/discord_ai_assistant .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/discord_ai_assistant .

RUN chmod +x ./discord_ai_assistant

ENTRYPOINT ["./discord_ai_assistant"]