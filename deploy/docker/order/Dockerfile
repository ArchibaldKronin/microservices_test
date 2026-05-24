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

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-order ./order/cmd/main.go


# ============================
# Stage 2: Final image
# ============================

FROM alpine:3.22

RUN apk add --no-cache curl

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser

COPY --from=builder /app/app-order .

COPY --from=builder /app/order/migrations ./migrations

EXPOSE 8080

ENTRYPOINT [ "./app-order" ]
