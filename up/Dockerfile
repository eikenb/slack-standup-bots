# Build stage - this container is discarded so size doesn't matter
FROM golang:alpine AS build-env
MAINTAINER jeikenberry@comscore.com
WORKDIR /opt/src/up
COPY . .
RUN apk --no-cache --update add git
ENV GOPATH=/opt
RUN go get .
RUN go build -ldflags='-w -s'

# Deployable stage - final image is based on this minimal image
FROM alpine:latest
RUN apk --no-cache --update add ca-certificates
WORKDIR /opt
COPY --from=build-env /opt/src/up/up ./
# runuprc: file containing exported environment variables needed by up
ARG runuprc=runuprc
ENV RUNUPRC=$runuprc
COPY $RUNUPRC runup ./
EXPOSE 8080
CMD ["/opt/runup"]
