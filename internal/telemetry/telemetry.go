package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	tele "gopkg.in/telebot.v4"
)

const instrumentationName = "github.com/ekucher/kbot"

var (
	tracer trace.Tracer

	updatesTotal    metric.Int64Counter
	errorsTotal     metric.Int64Counter
	handlerDuration metric.Float64Histogram
)

// Init initializes OpenTelemetry traces and metrics.
//
// OTLP exporter configuration is read from the standard
// OpenTelemetry environment variables.
func Init(ctx context.Context) (func(context.Context) error, error) {
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
	}

	metricExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		_ = traceExporter.Shutdown(ctx)
		return nil, fmt.Errorf("create OTLP metric exporter: %w", err)
	}

	serviceVersion := os.Getenv("APP_VERSION")
	if serviceVersion == "" {
		serviceVersion = "development"
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			attribute.String("service.name", "kbot"),
			attribute.String("service.version", serviceVersion),
		),
	)
	if err != nil {
		_ = traceExporter.Shutdown(ctx)
		_ = metricExporter.Shutdown(ctx)
		return nil, fmt.Errorf("create OpenTelemetry resource: %w", err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(
			traceExporter,
			sdktrace.WithBatchTimeout(2*time.Second),
		),
	)

	metricReader := sdkmetric.NewPeriodicReader(
		metricExporter,
		sdkmetric.WithInterval(5*time.Second),
	)

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(metricReader),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetMeterProvider(meterProvider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	tracer = traceProvider.Tracer(instrumentationName)
	meter := meterProvider.Meter(instrumentationName)

	updatesTotal, err = meter.Int64Counter(
		"kbot_updates",
		metric.WithDescription("Number of processed Telegram updates"),
	)
	if err != nil {
		_ = meterProvider.Shutdown(ctx)
		_ = traceProvider.Shutdown(ctx)
		return nil, fmt.Errorf("create updates counter: %w", err)
	}

	errorsTotal, err = meter.Int64Counter(
		"kbot_errors",
		metric.WithDescription("Number of Telegram handler errors"),
	)
	if err != nil {
		_ = meterProvider.Shutdown(ctx)
		_ = traceProvider.Shutdown(ctx)
		return nil, fmt.Errorf("create errors counter: %w", err)
	}

	handlerDuration, err = meter.Float64Histogram(
		"kbot_handler_duration",
		metric.WithDescription("Telegram handler execution duration"),
		metric.WithUnit("s"),
	)
	if err != nil {
		_ = meterProvider.Shutdown(ctx)
		_ = traceProvider.Shutdown(ctx)
		return nil, fmt.Errorf("create handler duration histogram: %w", err)
	}

	shutdown := func(ctx context.Context) error {
		return errors.Join(
			meterProvider.Shutdown(ctx),
			traceProvider.Shutdown(ctx),
		)
	}

	return shutdown, nil
}

// Handler wraps a Telebot handler with tracing, metrics and
// trace-correlated structured log fields.
func Handler(name string, next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		started := time.Now()

		handlerAttr := attribute.String("telegram.handler", name)

		ctx, span := tracer.Start(
			context.Background(),
			"telegram."+name,
			trace.WithAttributes(handlerAttr),
		)
		defer span.End()

		updatesTotal.Add(
			ctx,
			1,
			metric.WithAttributes(handlerAttr),
		)

		err := next(c)

		result := "ok"

		if err != nil {
			result = "error"

			errorsTotal.Add(
				ctx,
				1,
				metric.WithAttributes(handlerAttr),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		resultAttr := attribute.String("result", result)

		span.SetAttributes(resultAttr)

		duration := time.Since(started)

		handlerDuration.Record(
			ctx,
			duration.Seconds(),
			metric.WithAttributes(
				handlerAttr,
				resultAttr,
			),
		)

		spanContext := span.SpanContext()

		log.Printf(
			"event=telegram_update handler=%s result=%s duration_seconds=%.6f trace_id=%s span_id=%s",
			name,
			result,
			duration.Seconds(),
			spanContext.TraceID().String(),
			spanContext.SpanID().String(),
		)

		return err
	}
}
