# syntax=docker/dockerfile:1

ARG GO_VERSION=1.26
ARG DENO_VERSION=2.9.7

FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/ytdl-api ./cmd/api

FROM denoland/deno:bin-${DENO_VERSION} AS deno

FROM python:3.13-slim-trixie

ENV PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    PIP_ROOT_USER_ACTION=ignore

RUN apt-get update \
    && apt-get install -y --no-install-recommends ffmpeg ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ARG YTDLP_VERSION=
RUN pip install "yt-dlp[default]${YTDLP_VERSION:+==${YTDLP_VERSION}}"

COPY --from=deno /deno /usr/local/bin/deno
ENV DENO_NO_UPDATE_CHECK=1

RUN useradd --uid 1000 --create-home app \
    && mkdir -p /data/downloads /home/app/.cache \
    && chown app:app /data/downloads /home/app/.cache

COPY --from=build /out/ytdl-api /usr/local/bin/ytdl-api

USER app
WORKDIR /home/app
ENV API_PORT=8080 \
    DOWNLOAD_DIR=/data/downloads
EXPOSE 8080
ENTRYPOINT ["ytdl-api"]
