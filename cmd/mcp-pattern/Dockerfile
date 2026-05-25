# syntax=docker/dockerfile:1.7
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
    go build -trimpath -ldflags="-s -w" -o /out/mcp-pattern ./cmd/mcp-pattern

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.title="mcp-pattern"
LABEL org.opencontainers.image.licenses="Apache-2.0"
USER nonroot:nonroot
COPY --from=build /out/mcp-pattern /app/mcp-pattern
ENV TRANSPORT=stdio
ENV ADDR=:8087
EXPOSE 8087
ENTRYPOINT ["/app/mcp-pattern"]
