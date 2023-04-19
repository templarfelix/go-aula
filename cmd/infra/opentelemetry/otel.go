package opentelemetry

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"microservice/cmd/infra/config"
	"time"
)

func ProvideTracerProvider(logger *zap.SugaredLogger) *sdktrace.TracerProvider {
	logger.Info("Executing ProvideTracerProvider.")
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func ProvideMeterProvider(logger *zap.SugaredLogger) *metric.MeterProvider {
	logger.Info("Executing ProvideTracerProvider.")
	metricsExporter, err := stdoutmetric.New()
	if err != nil {
		panic(err)
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
	global.SetMeterProvider(provider)
	return provider
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
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return tracer.Shutdown(context.Background())
			},
		},
	)
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
				if err != nil {
					panic(err)
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return meter.Shutdown(context.Background())
			},
		},
	)
}
