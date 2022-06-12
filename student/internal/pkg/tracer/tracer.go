package tracer

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type Conf struct {
	Name string
	Env  string
	Ver  string
	Url  string
}

func NewConf(name, env, ver, url string) *Conf {
	return &Conf{
		Name: name,
		Env:  env,
		Ver:  ver,
		Url:  url,
	}
}

func (c *Conf) TracerProvider() (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.Url)),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(c.Name),
				attribute.String("env", c.Env),
				attribute.String("ver", c.Ver),
			)),
	)
	return tp, nil
}
