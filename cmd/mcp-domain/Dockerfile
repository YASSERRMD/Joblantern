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
    go build -trimpath -ldflags="-s -w" -o /out/mcp-domain ./cmd/mcp-domain

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.title="mcp-domain"
LABEL org.opencontainers.image.licenses="Apache-2.0"
USER nonroot:nonroot
COPY --from=build /out/mcp-domain /app/mcp-domain
ENV TRANSPORT=stdio
ENV ADDR=:8085
EXPOSE 8085
ENTRYPOINT ["/app/mcp-domain"]
