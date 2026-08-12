# Multi-stage prod image.
# Stage 1: build the React SPA.
FROM node:22-alpine AS frontend
WORKDIR /src/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: build the Go backend.
FROM golang:1.26-alpine AS backend
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /physiq ./cmd/server

# Stage 3: minimal runtime.
FROM alpine:3.21
RUN adduser -D -u 10001 physiq
WORKDIR /app
COPY --from=backend /physiq /usr/local/bin/physiq
COPY --from=frontend /src/frontend/dist /app/static
USER physiq
EXPOSE 8080
ENTRYPOINT ["physiq"]
