package main

import (
	"context"
	"fmt"
	"log"

	apm "go.stackify.com"
)

var (
	customKey = apm.Key("stackify.custom")
)

func main() {
	fmt.Println("Starting simple application.")

	stackifyAPM, err := apm.NewStackifyAPM(
		apm.WithApplicationName("Jayr GOLANG 11:22"),
		apm.WithEnvironmentName("Test"),
		apm.WithDebug(true),
	)
	defer stackifyAPM.Shutdown()

	if err != nil {
		log.Fatalf("failed to initialize stackifyapm: %v", err)
	}

	doSimpleSpan(stackifyAPM.Context, stackifyAPM.Tracer)
	doComplexSpan(stackifyAPM.Context, stackifyAPM.Tracer)

	fmt.Println("Application done.")
}

func doSimpleSpan(ctx context.Context, tracer apm.Tracer) {
	var err error = nil
	err = func(ctx context.Context) error {
		var span apm.Span
		ctx, span = tracer.Start(ctx, "custom.GoSampleClass.MethodCall")
		defer span.End()
		err = func(ctx context.Context) error {
			var span apm.Span
			ctx, span = tracer.Start(ctx, "custom.GoSampleClass2.MethodCall2")
			defer span.End()

			return nil
		}(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	}(ctx)
	if err != nil {
		panic(err)
	}
}

func doComplexSpan(ctx context.Context, tracer apm.Tracer) {
	var err error = nil
	err = func(ctx context.Context) error {
		var span apm.Span
		ctx, span = tracer.Start(ctx, "span1-0-0-0")
		defer span.End()

		err = func(ctx context.Context) error {
			var span apm.Span
			ctx, span = tracer.Start(ctx, "span1-1-0-0")
			defer span.End()

			err = func(ctx context.Context) error {
				var span apm.Span
				ctx, span = tracer.Start(ctx, "span1-1-1-0")
				defer span.End()

				err = func(ctx context.Context) error {
					var span apm.Span
					ctx, span = tracer.Start(ctx, "span1-1-1-1")
					defer span.End()

					return nil
				}(ctx)
				if err != nil {
					panic(err)
				}

				return nil
			}(ctx)
			if err != nil {
				panic(err)
			}

			return nil
		}(ctx)
		if err != nil {
			panic(err)
		}

		err = func(ctx context.Context) error {
			var span apm.Span
			ctx, span = tracer.Start(ctx, "span1-2-0-0")
			defer span.End()

			return nil
		}(ctx)
		if err != nil {
			panic(err)
		}

		err = func(ctx context.Context) error {
			var span apm.Span
			ctx, span = tracer.Start(ctx, "span1-3-0-0")
			defer span.End()
			err = func(ctx context.Context) error {
				var span apm.Span
				ctx, span = tracer.Start(ctx, "span1-3-1-0")
				defer span.End()

				return nil
			}(ctx)
			if err != nil {
				panic(err)
			}
			return nil
		}(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	}(ctx)
	if err != nil {
		panic(err)
	}
}
