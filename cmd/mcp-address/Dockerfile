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
    go build -trimpath -ldflags="-s -w" -o /out/mcp-address ./cmd/mcp-address

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.title="mcp-address"
LABEL org.opencontainers.image.licenses="Apache-2.0"
USER nonroot:nonroot
COPY --from=build /out/mcp-address /app/mcp-address
ENV TRANSPORT=stdio
ENV ADDR=:8082
EXPOSE 8082
ENTRYPOINT ["/app/mcp-address"]
