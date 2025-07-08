FROM golang:1.24 AS builder

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o tinypdf -trimpath -ldflags="-s -w" main.go

FROM alpine:latest

RUN apk add --no-cache \
    poppler-utils \
    qpdf \
    ghostscript

# Copy only the tinypdf script into the container at /app
COPY --from=builder /app/tinypdf /app/tinypdf

# Set permissions and move the script to path
RUN chmod +x /app/tinypdf && mv /app/tinypdf /usr/bin/

WORKDIR /app

# Run tinypdf when the container launches
ENTRYPOINT ["tinypdf"]