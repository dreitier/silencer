# syntax=docker/dockerfile:1

# builder image
FROM golang:1.18 AS builder

# if available, inject build args from GitHub Actions so that we have the current commit and tag
ARG GIT_COMMIT
ARG GIT_TAG

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
