FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY shared/ ./shared/

COPY retrival_mcp/ ./retrival_mcp/

WORKDIR /app/retrival_mcp/

RUN go build -v -ldflags="-w -s" -o /app/retrival_mcp .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/retrival_mcp .

RUN chmod +x ./retrival_mcp

ENTRYPOINT ["./retrival_mcp"]
