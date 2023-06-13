package tracerutil

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type TraceConf struct {
	Jaeger string `yaml:"jaeger"`
	Token  string `yaml:"token"`
}

func (c TraceConf) Trace(serviceName string) *sdktrace.TracerProvider {
	ctx := context.Background()
	fmt.Println(" - trace:", c.Jaeger)

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(c.Jaeger),
		otlptracegrpc.WithInsecure(),
	}

	res, err := resource.New(ctx,
		//设置 Token 值
		resource.WithAttributes(attribute.KeyValue{
			Key: "token", Value: attribute.StringValue(c.Token),
		}),
		//设置服务名
		resource.WithAttributes(attribute.KeyValue{
			Key: "service.name", Value: attribute.StringValue(serviceName),
		}),
	)
	if err != nil {
		panic(err)
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		panic(err)
	}

	//创建新的TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(10*time.Second)),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
