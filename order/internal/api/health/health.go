package health

import (
	"net/http"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"ok","service":"order-api"}`))
	if err != nil {
		http.Error(w, model.ErrInternal.Error(), http.StatusInternalServerError)
		logger.Error(r.Context(), "error healthcheck", zap.Error(err))
	}
}
