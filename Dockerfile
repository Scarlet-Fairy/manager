FROM golang:1.16-alpine AS base

MAINTAINER Michele Della Mea <michele.dellamea.arcanediver@gmail.com>

# Create appuser.
ARG USER=appuser
ARG UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /src

ENV CGO_ENABLED=0

COPY go.* ./
RUN go mod download
COPY . .

# ---------------------- #

FROM base AS build

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build \
    -o /out/manager \
    ./cmd/manager/main.go

# ---------------------- #

FROM scratch

COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=build /out/manager .

USER appuser:appuser

ENTRYPOINT ["/manager"]