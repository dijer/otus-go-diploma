FROM golang:1.22.2 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o previewer ./cmd/previewer

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/previewer /usr/local/bin/previewer

CMD ["previewer"]
