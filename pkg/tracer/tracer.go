package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type JeagerContext struct {
	Endpoint     string
	ServicesName string

	ctx context.Context
}

func (j JeagerContext) Open() (*sdktrace.TracerProvider, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(j.Endpoint)))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(j.ServicesName),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}

func (j *JeagerContext) SetConext(ctx context.Context) *JeagerContext {
	j.ctx = ctx
	return j
}

func (j *JeagerContext) WriteOtel(layer, spanName string) *JeagerContext {
	tracer := otel.Tracer(layer)
	ctx, span := tracer.Start(context.Background(), spanName)
	defer span.End()

	// fallback context
	j.SetConext(ctx)

	return j
}
