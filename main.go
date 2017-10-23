package main

import (
	"context"
	"os"

	lightstep "github.com/lightstep/lightstep-tracer-go"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

var (
	envToken           = os.Getenv("LIGHTSTEP_ACCESS_TOKEN")
	envCustomCollector = os.Getenv("BACKPLANE_LIGHTSTEP_ENDPOINT")
)

func main() {
	if envToken == "" {
		panic("LIGHTSTEP_ACCESS_TOKEN not set.")
	}

	lsOpts := lightstep.Options{
		AccessToken: envToken,
		UseGRPC:     true,
	}

	if envCustomCollector != "" {
		lsOpts.Collector = lightstep.Endpoint{
			Host:      envCustomCollector,
			Port:      443,
			Plaintext: false,
		}
	}

	opentracing.InitGlobalTracer(lightstep.NewTracer(lsOpts))

	hello(context.Background())

	err := lightstep.FlushLightStepTracer(opentracing.GlobalTracer())
	if err != nil {
		panic(err)
	}
}

func hello(ctx context.Context) {
	sp, _ := opentracing.StartSpanFromContext(ctx, "test span")
	defer sp.Finish()

	sp.LogEvent("HELLO LIGHTSTEP. ARE YOU THERE!?")
	sp.LogFields(log.String("hello_lightstep", "love_xoxo"), log.Object("testSpan", sp))
}
