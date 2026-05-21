package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName = "assembly-service"
)

var meter = otel.Meter(serviceName)

// assemblyDurationSeconds Гистограмма длительности сборки
var assemblyDurationSeconds metric.Float64Histogram

func InitMetrics() error {
	var err error

	assemblyDurationSeconds, err = meter.Float64Histogram(
		"assembly_duration_seconds",
		metric.WithDescription("Duration of assebly order"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(
			0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0, 10.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}

func AppendAssemblyDurationSecondsMetric(ctx context.Context, duration time.Duration) {
	durationSeconds := duration.Seconds()
	assemblyDurationSeconds.Record(ctx, durationSeconds)
}
