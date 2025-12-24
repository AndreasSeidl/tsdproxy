
# Stage 1: Build web assets
FROM oven/bun:1 AS web-builder
WORKDIR /app/web
COPY web/package.json web/bun.lockb* ./
RUN bun install
COPY web/ ./
RUN bun run build

# Stage 2: Build Go application
FROM golang:1.25 AS builder
RUN apk add --no-cache ca-certificates && update-ca-certificates 2>/dev/null || true

# Define o diretório de trabalho
WORKDIR /app

# Copia o código fonte para o container
COPY . .

# Copy the built web assets from stage 1
COPY --from=web-builder /app/web/dist ./web/dist

# Install templ for template generation (match the version in go.mod)
RUN go install github.com/a-h/templ/cmd/templ@v0.3.865

# Generate templ templates
RUN templ generate

# Compila a aplicação Go
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /tsdproxyd ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /healthcheck ./cmd/healthcheck/main.go

# Stage 3: Final image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /tsdproxyd /tsdproxyd
COPY --from=builder /healthcheck /healthcheck

ENTRYPOINT ["/tsdproxyd"]

EXPOSE 8080
HEALTHCHECK CMD [ "/healthcheck" ]
