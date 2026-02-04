package processor

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// mockExporter is a test implementation of ItemExporter.
type mockExporter[T any] struct {
	exportedItems []*T
	exportCount   atomic.Int64
	exportErr     error
	exportDelay   time.Duration
}

func (m *mockExporter[T]) ExportItems(ctx context.Context, items []*T) error {
	if m.exportDelay > 0 {
		select {
		case <-time.After(m.exportDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	m.exportCount.Add(int64(len(items)))
	m.exportedItems = append(m.exportedItems, items...)

	return m.exportErr
}

func (m *mockExporter[T]) Shutdown(_ context.Context) error {
	return nil
}

func TestBatchItemProcessor_Basic(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	exporter := &mockExporter[string]{}

	proc, err := NewBatchItemProcessor[string](
		exporter,
		"test",
		log,
		WithMaxQueueSize(1000),
		WithMaxExportBatchSize(10),
		WithBatchTimeout(100*time.Millisecond),
		WithWorkers(1),
	)
	if err != nil {
		t.Fatalf("failed to create processor: %v", err)
	}

	ctx := context.Background()
	proc.Start(ctx)

	// Write some items.
	items := make([]*string, 5)
	for i := range items {
		s := "item"
		items[i] = &s
	}

	if err := proc.Write(ctx, items); err != nil {
		t.Fatalf("failed to write items: %v", err)
	}

	// Wait for batch timeout to trigger export.
	time.Sleep(200 * time.Millisecond)

	// Shutdown.
	if err := proc.Shutdown(ctx); err != nil {
		t.Fatalf("failed to shutdown: %v", err)
	}

	// Verify items were exported.
	if exporter.exportCount.Load() != 5 {
		t.Errorf("expected 5 items exported, got %d", exporter.exportCount.Load())
	}
}

func TestBatchItemProcessor_BatchSize(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	exporter := &mockExporter[int]{}

	proc, err := NewBatchItemProcessor[int](
		exporter,
		"test",
		log,
		WithMaxQueueSize(1000),
		WithMaxExportBatchSize(5),
		WithBatchTimeout(10*time.Second), // Long timeout to ensure batch size triggers.
		WithWorkers(1),
	)
	if err != nil {
		t.Fatalf("failed to create processor: %v", err)
	}

	ctx := context.Background()
	proc.Start(ctx)

	// Write exactly batch size items.
	items := make([]*int, 5)
	for i := range items {
		val := i
		items[i] = &val
	}

	if err := proc.Write(ctx, items); err != nil {
		t.Fatalf("failed to write items: %v", err)
	}

	// Wait briefly for processing.
	time.Sleep(100 * time.Millisecond)

	// Should have exported due to batch size.
	if exporter.exportCount.Load() != 5 {
		t.Errorf("expected 5 items exported, got %d", exporter.exportCount.Load())
	}

	if err := proc.Shutdown(ctx); err != nil {
		t.Fatalf("failed to shutdown: %v", err)
	}
}

func TestBatchItemProcessor_OptionsValidation(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	exporter := &mockExporter[string]{}

	// Test invalid: batch size > queue size.
	_, err := NewBatchItemProcessor[string](
		exporter,
		"test",
		log,
		WithMaxQueueSize(10),
		WithMaxExportBatchSize(20),
	)

	if err == nil {
		t.Error("expected error for batch size > queue size")
	}

	// Test invalid: zero workers.
	_, err = NewBatchItemProcessor[string](
		exporter,
		"test",
		log,
		WithWorkers(0),
	)

	if err == nil {
		t.Error("expected error for zero workers")
	}
}
