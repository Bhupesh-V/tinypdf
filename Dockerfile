FROM golang:1.24 AS builder

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o tinypdf -trimpath -ldflags="-s -w" main.go

FROM alpine:latest AS tinypdf-ops

RUN apk add --no-cache \
    poppler-utils \
    qpdf \
    ghostscript

# Copy only the tinypdf script into the container at /app
COPY --from=builder /app/tinypdf .

# Set permissions and move the script to path
RUN chmod +x tinypdf && mv tinypdf /usr/bin/

LABEL description="🤏🏽 Reduce PDF file size"

# Create tmp directory for temporary files
WORKDIR /tmp

# Set environment variable for temporary directory
ENV TMPDIR=/tmp

WORKDIR /app

# Run tinypdf when the container launches
ENTRYPOINT ["tinypdf"]