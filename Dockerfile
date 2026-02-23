FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o hello-world-api .

FROM scratch
COPY --from=builder /app/hello-world-api /hello-world-api

EXPOSE 8080
ENTRYPOINT ["/hello-world-api"]
