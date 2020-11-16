# Stackify Go APM

## Installation Guide

### Standalone setup

1. Install **Stackify Linux Agent**.

2. Check that your setup meets our system requirements.
    - Go Version 1.15

3. Install the Stackify Go APM agent using `go get`
    - Add stackify go apm to your `go.mod`

        ```
        require (
            github.com/stackify/stackify-go-apm vx.x.x
            ...
        )
        ```

    - Install stackify go apm

        ```
        $ go get github.com/stackify/stackify-go-apm
        ```

4. Update and insert the apm settings to your application.

    ```
    package main

    import (
        "context"
        "log"

        "github.com/stackify/stackify-go-apm"
    )

    func main() {
        stackifyAPM, err := apm.NewStackifyAPM()
        if err != nil {
            log.Fatalf("failed to initialize stackifyapm: %v", err)
        }
        defer stackifyAPM.Shutdown()

        err = func(ctx context.Context) error {
            var span apm.Span
            ctx, span = stackifyAPM.Tracer.Start(ctx, "span1")
            defer span.End()

            return nil
        }(stackifyAPM.Context)
        if err != nil {
            panic(err)
        }

    }
    ```

5. Customize **Application Name** and **Environment** configuration.

    ```
    stackifyAPM, err := apm.NewStackifyAPM(
        apm.WithApplicationName("Go Application"),
        apm.WithEnvironmentName("Production"),
    )
    ```
