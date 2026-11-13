FROM node:22.21.0-alpine AS frontend
WORKDIR /app
RUN corepack enable
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend ./frontend
COPY vite.config.ts tsconfig.json ./
RUN pnpm build

FROM golang:1.26.11-alpine AS backend
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum* ./
RUN go mod download
COPY backend ./
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.22
WORKDIR /app
COPY --from=backend /server ./server
COPY --from=frontend /app/frontend/dist ./frontend/dist
ENV FRONTEND_DIR=/app/frontend/dist
EXPOSE 8080
USER 65532:65532
CMD ["./server"]
