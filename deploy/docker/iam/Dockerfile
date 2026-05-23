# ============================
# Stage 1: Build stage
# ============================
FROM golang:1.26.2-alpine3.22 AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.work go.work.sum ./

COPY assembly/go.mod assembly/go.sum ./assembly/
COPY iam/go.mod iam/go.sum ./iam/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY notification/go.mod notification/go.sum ./notification/
COPY order/go.mod order/go.sum ./order/
COPY payment/go.mod payment/go.sum ./payment/
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/

RUN go mod download

COPY assembly ./assembly
COPY iam ./iam
COPY inventory ./inventory
COPY notification ./notification
COPY order ./order
COPY payment ./payment
COPY platform ./platform
COPY shared ./shared

# COPY bin/grpc-health-probe ./
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe

RUN chmod +x grpc-health-probe

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-iam ./iam/cmd/main.go

# ============================
# Stage 2: Final image
# ============================

FROM alpine:3.22

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser

COPY --from=builder /app/app-iam .

COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

COPY --from=builder /app/iam/migrations ./migrations

HEALTHCHECK --interval=10s --timeout=3s --retries=3 CMD grpc-health-probe -addr=localhost:50053 || exit 1

EXPOSE 50053

ENTRYPOINT [ "./app-iam" ]
