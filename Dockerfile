FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o /ms365-mcp-server ./cmd/ms365-mcp-server

FROM alpine:3.20
COPY --from=build /ms365-mcp-server /usr/local/bin/ms365-mcp-server
ENTRYPOINT ["/usr/local/bin/ms365-mcp-server"]
