# go-batch-processor

A generic Go library for batching items and exporting them efficiently.

## Installation

```bash
go get github.com/ethpandaops/go-batch-processor
```

## Usage

```go
package main

import (
    "context"
    "github.com/ethpandaops/go-batch-processor/processor"
    "github.com/sirupsen/logrus"
)

// Implement the ItemExporter interface
type MyExporter struct{}

func (e *MyExporter) ExportItems(ctx context.Context, items []*MyItem) error {
    // Send items to your destination (API, database, etc.)
    return nil
}

func (e *MyExporter) Shutdown(ctx context.Context) error {
    return nil
}

func main() {
    log := logrus.New()
    exporter := &MyExporter{}

    proc, _ := processor.NewBatchItemProcessor[MyItem](
        exporter,
        "my-processor",
        log,
        processor.WithMaxExportBatchSize(100),
        processor.WithBatchTimeout(5*time.Second),
        processor.WithWorkers(3),
    )

    ctx := context.Background()
    proc.Start(ctx)
    defer proc.Shutdown(ctx)

    // Write items - they'll be batched and exported automatically
    proc.Write(ctx, []*MyItem{{}, {}})
}
```

## Configuration

| Option | Default | Description |
|--------|---------|-------------|
| `WithMaxQueueSize` | 51,200 | Maximum items to buffer |
| `WithMaxExportBatchSize` | 512 | Items per batch |
| `WithBatchTimeout` | 5s | Max wait before sending partial batch |
| `WithExportTimeout` | 30s | Timeout for export operations |
| `WithWorkers` | 5 | Concurrent export workers |
| `WithShippingMethod` | Async | `ShippingMethodAsync` or `ShippingMethodSync` |

## Features

- Generic type support (`[T any]`)
- Async and sync shipping modes
- Configurable batch size and timeout triggers
- Worker pool for concurrent exports
- Built-in Prometheus metrics
- Graceful shutdown with queue draining

## License

See [LICENSE](LICENSE).
