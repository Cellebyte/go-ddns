# This Dockerfile is used to build the image available on DockerHub
FROM docker.io/golang:1.25-alpine AS build

RUN apk --no-cache add ca-certificates tzdata

# Add everything
ADD . /usr/src/go-ddns

RUN cd /usr/src/go-ddns && \
    CGO_ENABLED=0 go build -o go-ddns .

FROM scratch
LABEL org.opencontainers.image.source=https://github.com/cellebyte/go-ddns
COPY --from=build /usr/src/go-ddns/go-ddns /usr/bin/go-ddns
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=UTC

WORKDIR /usr/bin

ENTRYPOINT ["go-ddns"]