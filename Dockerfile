FROM vault:latest as build
RUN apk add --update alpine-sdk
RUN apk update && apk add go git
WORKDIR /app
ENV GOPATH /app
ADD . /app/src
RUN GO111MODULE=on go get github.com/immutability-io/vault-ethereum
RUN GO111MODULE=on CGO_ENABLED=1 GOOS=linux go install -a github.com/immutability-io/vault-ethereum

FROM vault:latest
WORKDIR /app
RUN cd /app
COPY --from=build /app/bin/vault-ethereum /app/bin/vault-ethereum
# Prove the binary is now an executable
CMD /app/bin/vault-ethereum