package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func ReadyCheck(pool *pgxpool.Pool, clientsConnections ...*grpc.ClientConn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if pool == nil {
			http.Error(w, "db not configured", http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		errG, ctx := errgroup.WithContext(ctx)

		errG.Go(func() error {
			ctxDB, cancelDB := context.WithTimeout(ctx, 3*time.Second)
			defer cancelDB()

			if err := pool.Ping(ctxDB); err != nil {
				return fmt.Errorf("db not ready: %w", err)
			}
			return nil
		})

		for i, clientCon := range clientsConnections {
			i := i
			clientCon := clientCon
			errG.Go(func() error {
				if clientCon == nil {
					return fmt.Errorf("client %d is nil", i)
				}

				hc := grpc_health_v1.NewHealthClient(clientCon)

				ctxCl, cancelCl := context.WithTimeout(ctx, 3*time.Second)
				defer cancelCl()

				resp, err := hc.Check(
					ctxCl,
					&grpc_health_v1.HealthCheckRequest{},
				)
				if err != nil {
					return fmt.Errorf("client %d error: %w", i, err)
				}
				if resp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
					return fmt.Errorf("client %d not serving: %s", i, resp.GetStatus())
				}
				return nil
			})
		}

		err := errG.Wait()
		if err != nil {
			http.Error(w, fmt.Sprintf("service is not ready: %v", err), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			logger.Error(r.Context(), "error readycheck", zap.Error(err))
		}
	}
}
