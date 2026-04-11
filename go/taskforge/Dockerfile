
# Build Stage
FROM golang:1.23 AS builder 

WORKDIR /app 

# Cache dependecies first 
COPY go.mod go.sum ./ 
RUN go mod download 

# Copy source 
COPY . .

# Build the API binary 
RUN go build -o taskforge-api ./cmd/taskforge-api 

# Runtime stage 
FROM debian:bookworm-slim

WORKDIR /app 

COPY --from=builder /app/taskforge-api /app/taskforge-api 

EXPOSE 8080 

CMD ["./taskforge-api"]
