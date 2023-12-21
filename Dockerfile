# syntax=docker/dockerfile:1

ARG BASE_IMAGE=golang:1.21

FROM $BASE_IMAGE AS build-stage

ARG SUB_PROJECT

WORKDIR /srv

COPY . .

WORKDIR /srv/$SUB_PROJECT

RUN --mount=type=secret,id=netrc,dst=/root/.netrc CGO_ENABLED=0 go build -o bin/appbin

FROM gcr.io/distroless/static-debian11:nonroot

ARG SUB_PROJECT

WORKDIR /srv

COPY --from=build-stage /srv/${SUB_PROJECT}/bin/appbin ./appbin

ENTRYPOINT ["/srv/appbin"]
