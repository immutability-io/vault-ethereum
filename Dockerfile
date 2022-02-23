FROM golang:1.17-alpine as builder
RUN apk add --update alpine-sdk
RUN apk update && apk add git openssh gcc musl-dev linux-headers

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY  / .
RUN mkdir -p /build/bin \
    && CGO_ENABLED=1 GOOS=linux go build -a -v -i -o /build/bin/vault-ethereum . \
    && sha256sum -b /build/bin/vault-ethereum > /build/bin/SHA256SUMS

FROM vault:latest
ARG always_upgrade
RUN echo ${always_upgrade} > /dev/null && apk update && apk upgrade
RUN apk add bash openssl jq

USER vault
WORKDIR /app
RUN mkdir -p /home/vault/plugins

COPY --from=builder /build/bin/vault-ethereum /home/vault/plugins/vault-ethereum
COPY --from=builder /build/bin/SHA256SUMS /home/vault/plugins/SHA256SUMS
RUN ls -la /home/vault/plugins
HEALTHCHECK CMD nc -zv 127.0.0.1 9200 || exit 1
