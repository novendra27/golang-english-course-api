# ==========================================
# Stage 1: Build Stage (Kompilasi Binary Go)
# ==========================================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git dan dependency tools jika diperlukan
RUN apk add --no-cache git tzdata

# Cache layer dependensi Go
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Kompilasi binary statis (CGO_ENABLED=0 agar tidak bergantung pada C library sistem)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server main.go

# ==========================================
# Stage 2: Final Lightweight Runtime Stage
# ==========================================
FROM alpine:3.20

WORKDIR /app

# Install timezone data & certificates
RUN apk --no-cache add ca-certificates tzdata

# Copy binary yang sudah dikompilasi dari stage builder
COPY --from=builder /app/server /app/server

# Expose port aplikasi
EXPOSE 8080

# Jalankan server
CMD ["/app/server"]
