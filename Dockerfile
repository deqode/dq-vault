# Stage 1 (to create a "build" image)
FROM golang:1.10.1 AS source

COPY . /go/src/gitlab.com/arout/Vault/
WORKDIR /go/src/gitlab.com/arout/Vault/

RUN go build

# Stage 2 (to create a vault conatiner with executable)
FROM vault:latest

# Make new directory for plugins
RUN mkdir /vault/plugins && \
  apk --no-cache add ca-certificates wget && \
  wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
  wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.28-r0/glibc-2.28-r0.apk && \
  apk add glibc-2.28-r0.apk

# Copy executable from source to vault
COPY --from=source /go/src/gitlab.com/arout/Vault/Vault /vault/plugins/Vault

# TODO: add make run