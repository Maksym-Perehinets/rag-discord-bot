FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY . .

RUN go build -v -ldflags="-w -s" -o /app/message_parsing_executable .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/message_parsing_executable ./message_parsing

RUN chmod +x ./message_parsing

ENTRYPOINT ["./message_parsing"]