FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd

FROM surnet/alpine-wkhtmltopdf:3.20.0-0.12.6-full

WORKDIR /app

COPY --from=builder /app/server .

COPY templates ./templates
COPY migrations ./migrations


EXPOSE 8080

ENTRYPOINT ["./server"]
