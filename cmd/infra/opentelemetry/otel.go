package opentelemetry

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"microservice/cmd/infra/config"
)

func ProvideTracerProvider(logger *zap.SugaredLogger) (*sdktrace.TracerProvider, error) {
	logger.Info("Executing ProvideTracerProvider.")
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		logger.Errorw("failed to create trace exporter", "error", err)
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func ProvideMeterProvider(logger *zap.SugaredLogger) (*metric.MeterProvider, error) {
	logger.Info("Executing ProvideMeterProvider.")
	metricsExporter, err := stdoutmetric.New()
	if err != nil {
		logger.Errorw("failed to create metrics exporter", "error", err)
		return nil, err
	}
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)
	reader := metric.NewPeriodicReader(metricsExporter)
	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)
	otel.SetMeterProvider(provider)
	return provider, nil
}

func RegisterHooks(
	lifecycle fx.Lifecycle,
	logger *zap.SugaredLogger,
	config config.Config,
	echoInstance *echo.Echo,
	tracer *sdktrace.TracerProvider,
	meter *metric.MeterProvider,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				echoInstance.Use(otelecho.Middleware(config.App))
				if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
					logger.Errorw("failed to start runtime instrumentation", "error", err)
					return err
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("Shutting down OpenTelemetry providers...")
				if err := tracer.Shutdown(ctx); err != nil {
					logger.Errorw("failed to shutdown tracer provider", "error", err)
				}
				if err := meter.Shutdown(ctx); err != nil {
					logger.Errorw("failed to shutdown meter provider", "error", err)
				}
				return nil
			},
		},
	)
}
