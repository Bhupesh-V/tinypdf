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

# Third stage: Copy only necessary binaries and their dependencies
FROM scratch

LABEL description="🤏🏽 Reduce PDF file size"
LABEL maintainer="Bhupesh Varshney <varshneybhupesh@gmail.com>"

# Create tmp directory for temporary files
WORKDIR /tmp

# Set environment variable for temporary directory
ENV TMPDIR=/tmp

# Copy the main binary
COPY --from=tinypdf-ops /usr/bin/tinypdf /usr/bin/

# Copy essential system libraries
COPY --from=tinypdf-ops /lib/ld-musl-*.so.1 /lib/
COPY --from=tinypdf-ops /lib/libc.musl-*.so.1 /lib/
COPY --from=tinypdf-ops /usr/lib/libgcc_s.so.1 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libstdc++.so.6 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libz.so* /usr/lib/
COPY --from=tinypdf-ops /lib/libz.so* /lib/

# Copy qpdf binary
COPY --from=tinypdf-ops /usr/bin/qpdf /usr/bin/
# Copy qpdf dependencies
COPY --from=tinypdf-ops /usr/lib/libqpdf.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libcrypto.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/liblzma.so* /usr/lib/

# Copy poppler-utils binaries
# COPY --from=tinypdf-ops /usr/bin/pdfinfo /usr/bin/
# COPY --from=tinypdf-ops /usr/bin/pdftoppm /usr/bin/
# COPY --from=tinypdf-ops /usr/bin/pdftotext /usr/bin/
# COPY --from=tinypdf-ops /usr/bin/pdftops /usr/bin/
# COPY --from=tinypdf-ops /usr/bin/pdfunite /usr/bin/
# COPY --from=tinypdf-ops /usr/bin/pdfseparate /usr/bin/

# Copy pdftocairo dependencies
COPY --from=tinypdf-ops /usr/bin/pdftocairo /usr/bin/
COPY --from=tinypdf-ops /usr/lib/libcairo.so.2 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libpoppler.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libXext.so.6 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libXrender.so.1 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libxcb-render.so.0 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libxcb-shm.so.0 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libpixman-1.so.0 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libsmime* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libnss3.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libplc4.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libnspr4.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libnssutil3.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libplds4.so* /usr/lib/

# Copy ghostscript binary
COPY --from=tinypdf-ops /usr/bin/gs /usr/bin/
# Copy ghostscript data files
COPY --from=tinypdf-ops /usr/share/ghostscript/ /usr/share/ghostscript
# Copy ghostscript dependencies
COPY --from=tinypdf-ops /usr/lib/libXt.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libX11.so.6 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libtiff.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libcups.so.2 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libpng16.so.16 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libjbig2dec.so.0 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libjpeg.so.8 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/liblcms2.so.2 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libfontconfig.so.1 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libfreetype.so.6 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libopenjp2.so.7 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libSM.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libICE.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libxcb.so.1 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libzstd.so.1 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libwebp.so.7 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libavahi-common.so.3 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libavahi-client.so.3 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libgnutls.so.30 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libexpat.so.1 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libbz2.so.1 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libbrotlidec.so.1 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libuuid.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libXau.so.6 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libXdmcp.so.6 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libsharpyuv.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libintl.so.8 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libdbus-1.so.3 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libp11-kit.so.0 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libidn2.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libunistring.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libtasn1.so.6 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libhogweed.so.6 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libnettle.so.8 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libgmp.so.10 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libbrotlicommon.so.1 /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libbsd.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libffi.so* /usr/lib/
COPY --from=tinypdf-ops /usr/lib/libmd.so* /usr/lib/

WORKDIR /app

# Run tinypdf when the container launches
ENTRYPOINT ["tinypdf"]