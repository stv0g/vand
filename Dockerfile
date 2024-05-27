# SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.22-alpine AS backend-builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -o vand ./cmd

# FROM node:17 AS frontend-builder

# ENV NODE_ENV=production

# WORKDIR /app

# COPY frontend/package.json .
# COPY frontend/package-lock.json* .

# RUN npm install --production

# COPY frontend/ .

# RUN npm run build

FROM alpine:3.20

RUN apk update && apk add ca-certificates curl && rm -rf /var/cache/apk/*

# COPY --from=frontend-builder /app/dist/ /dist/
COPY --from=backend-builder /app/vand /
COPY --from=backend-builder /app/etc/vand.yaml /

ENV GIN_MODE=release

EXPOSE 8080/tcp

HEALTHCHECK --interval=30s --timeout=30s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/healthz

ENTRYPOINT [ "/vand" ]
