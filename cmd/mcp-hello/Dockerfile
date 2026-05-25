# syntax=docker/dockerfile:1.7
#
# joblantern.hello — reference MCP server.
# Multi-stage: build with golang-alpine, ship on distroless nonroot.

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
    -o /out/mcp-hello ./cmd/mcp-hello

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.title="mcp-hello"
LABEL org.opencontainers.image.source="https://github.com/yasserrmd/joblantern"
LABEL org.opencontainers.image.licenses="Apache-2.0"
USER nonroot:nonroot
COPY --from=build /out/mcp-hello /app/mcp-hello
ENV TRANSPORT=stdio
ENV ADDR=:8081
EXPOSE 8081
ENTRYPOINT ["/app/mcp-hello"]
