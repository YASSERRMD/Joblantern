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
    go build -trimpath -ldflags="-s -w" -o /out/mcp-eu-vies ./cmd/mcp-eu-vies

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.licenses="Apache-2.0"
USER nonroot:nonroot
COPY --from=build /out/mcp-eu-vies /app/mcp-eu-vies
ENV TRANSPORT=stdio
ENV ADDR=:8092
EXPOSE 8092
ENTRYPOINT ["/app/mcp-eu-vies"]
