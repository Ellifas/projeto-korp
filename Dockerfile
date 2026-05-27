# Stage 1: build
FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/http-server-projeto-korp ./cmd/server

# Stage 2: runtime
FROM alpine:3.19

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=build /app/http-server-projeto-korp /app/http-server-projeto-korp

EXPOSE 8080

USER appuser

CMD ["/app/http-server-projeto-korp"]
