# syntax=docker/dockerfile:1

# builder image
FROM golang:1.18 AS builder
WORKDIR /build
COPY . ./
RUN make

# target image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/silencer .

RUN mkdir /etc/silencer \
    && echo "---"  > /etc/silencer/config.yaml \
    && adduser -D -u 1000 silencer

USER silencer
ENTRYPOINT ["/app/silencer"]
