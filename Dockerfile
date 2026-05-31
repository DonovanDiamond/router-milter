
FROM golang:1.26-alpine AS builder
ARG VERSION
ARG COMMIT
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" -o router-milter .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/router-milter .
ENTRYPOINT ["./router-milter"]
CMD ["-config", "/app/config.yaml"]
