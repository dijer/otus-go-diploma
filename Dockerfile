FROM golang:1.22.2 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o previewer ./cmd/previewer

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/previewer /usr/local/bin/previewer

CMD ["previewer"]
