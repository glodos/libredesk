FROM golang:1.27-alpine AS builder

RUN apk --no-cache add make git nodejs npm grep
RUN npm install -g pnpm@11.28.8

WORKDIR /src
COPY . .

# This fork is based on the latest upstream LibreDesk release.
# Makefile uses LIBREDESK_VERSION before falling back to git tags / v0.0.0.
ARG LIBREDESK_VERSION=v2.9.2
ENV LIBREDESK_VERSION=${LIBREDESK_VERSION}

# Build frontend, backend, and embed static assets into the single LibreDesk binary.
RUN make build

FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /libredesk

COPY --from=builder /src/libredesk .
COPY config.sample.toml config.toml

EXPOSE 9000

CMD ["./libredesk"]
