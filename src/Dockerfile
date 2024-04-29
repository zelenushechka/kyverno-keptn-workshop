# Stage 1: Build the Go application
FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Stage 2: Create a minimal image with the statically linked binary
FROM scratch

COPY --from=builder /app/server /app/server
COPY index.html /app/index.html
COPY kcd.png /app/kcd.png

EXPOSE 8080

ENTRYPOINT ["/app/server"]