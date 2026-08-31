# build
FROM golang:alpine AS server_builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct 
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ohurlshortener

# deploy
FROM alpine:latest 
RUN apk --no-cache add ca-certificates tzdata wget
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=server_builder /app/ohurlshortener . 
COPY --from=server_builder /app/.env .
COPY --from=server_builder /app/config_docker.ini ./config.ini
COPY --from=server_builder /app/templates ./templates
EXPOSE 9091 9092
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9091/healthz || exit 1
CMD ["./ohurlshortener"]
