# syntax=docker/dockerfile:1.7
#
# joblantern main app — multi-stage build.
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/joblantern ./cmd/joblantern

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.title="joblantern"
LABEL org.opencontainers.image.source="https://github.com/yasserrmd/joblantern"
LABEL org.opencontainers.image.licenses="Apache-2.0"
USER nonroot:nonroot
WORKDIR /app
COPY --from=build /out/joblantern /app/joblantern
ENV JOBLANTERN_ADDR=:8080
EXPOSE 8080
ENTRYPOINT ["/app/joblantern"]
