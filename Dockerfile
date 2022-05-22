FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git dep openssh-client

COPY .  ${GOPATH}/src/silencer
RUN cd ${GOPATH}/src/silencer \
    && dep ensure -v \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o silencer \
    && strip silencer


# Alpine doesn't work, archives aren't created. Probably due to musl libc
FROM fedora:31
WORKDIR /usr/local/bin

RUN mkdir /etc/silencer \
    && echo "---"  > /etc/silencer/config.yaml \
    && useradd -u 1000 silencer

USER silencer

COPY --from=build-env /go/src/silencer/silencer .

ENTRYPOINT ["/usr/local/bin/silencer"]