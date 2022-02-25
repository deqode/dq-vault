# Stage 1 (to create a "build" image)
FROM golang:1.10.1 AS source

RUN curl https://glide.sh/get | sh

COPY . /go/src/github.com/deqode/dq-vault/
WORKDIR /go/src/github.com/deqode/dq-vault/

#RUN glide install
RUN go build

# Stage 2 (to create a vault conatiner with executable)
FROM vault:latest

# Make new directory for plugins
RUN mkdir /vault/plugins

RUN apk --no-cache add ca-certificates wget make
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub
RUN wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.28-r0/glibc-2.28-r0.apk
RUN apk add glibc-2.28-r0.apk


# Copy executable from source to vault
COPY --from=source /go/src/github.com/deqode/dq-vault/dq-vault /vault/plugins/vault_plugin
COPY ./Makefile .
COPY ./setup/config.hcl /vault/config/config.hcl

# TODO: add make run
CMD ["make", "run"]
