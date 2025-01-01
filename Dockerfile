FROM golang:1.26.5 AS builder

ARG ACTX0_VERSION=0.1.0
ARG ACTX0_COMMIT=none
ARG ACTX0_BUILD_DATE=unknown
ARG ACTX0_BUILT_BY=docker

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${ACTX0_VERSION} -X main.commit=${ACTX0_COMMIT} -X main.date=${ACTX0_BUILD_DATE} -X main.builtBy=${ACTX0_BUILT_BY}" \
    -o /out/actx0 .

RUN mkdir -p /out/configs /out/var/logs && \
    cp config.dist.yml /out/configs/config.dist.yml

FROM gcr.io/distroless/static-debian13:nonroot

ARG ACTX0_VERSION=0.1.0
ARG ACTX0_COMMIT=none
ARG ACTX0_BUILD_DATE=unknown
ARG ACTX0_BUILT_BY=docker

LABEL org.opencontainers.image.title="actx0" \
      org.opencontainers.image.description="actx0 platform" \
      org.opencontainers.image.url="https://github.com/clivern/actx0" \
      org.opencontainers.image.source="https://github.com/clivern/actx0" \
      org.opencontainers.image.version="${ACTX0_VERSION}" \
      org.opencontainers.image.revision="${ACTX0_COMMIT}" \
      org.opencontainers.image.created="${ACTX0_BUILD_DATE}" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.vendor="Clivern" \
      org.opencontainers.image.authors="Clivern" \
      built.by="${ACTX0_BUILT_BY}"

WORKDIR /app

COPY --from=builder --chown=65532:65532 /out/actx0 /app/actx0
COPY --from=builder --chown=65532:65532 /out/configs /app/configs
COPY --from=builder --chown=65532:65532 /out/var /app/var

EXPOSE 8080

VOLUME ["/app/configs", "/app/var"]

USER 65532:65532

ENTRYPOINT ["/app/actx0"]

CMD ["server", "-c", "/app/configs/config.dist.yml"]
