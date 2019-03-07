package jaeger

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

// 初始化变量
var (
	err    error
	Tracer opentracing.Tracer
	Closer io.Closer
)

// SpanContextKey is context span id
const SpanContextKey = "otspan"

// Init returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func Init(serviceName string) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "probabilistic",
			Param: 0.3,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	Tracer, Closer, err = cfg.New(serviceName, config.Logger(jaeger.NullLogger))
	if err != nil {
		fmt.Printf("ERROR: cannot init Jaeger: %s\n", err.Error())
	}
	opentracing.SetGlobalTracer(Tracer)
}

// StartSpanFromHeader is start span from the request header.
func StartSpanFromHeader(header *http.Header, operationName string) (span opentracing.Span, ctx context.Context) {
	spanCtx, _ := Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(*header))
	span, ctx = opentracing.StartSpanFromContext(context.Background(), operationName, ext.RPCServerOption(spanCtx))
	return span, ctx
}

// StartSpanFromContext is start span from the context.
func StartSpanFromContext(ctx context.Context, operationName string) (span opentracing.Span, childCtx context.Context) {
	span, childCtx = opentracing.StartSpanFromContext(ctx, operationName)
	return span, childCtx
}

// StartSpanFromParentSpan is start span from the parent span.
func StartSpanFromParentSpan(parentSpan opentracing.Span, operationName string) (span opentracing.Span) {
	span = opentracing.StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
	return span
}

// GetSpanFormContext is get span from the echo context.
func GetSpanFormContext(ctx echo.Context) (span opentracing.Span) {
	span = ctx.Get(SpanContextKey).(opentracing.Span)
	return span
}
