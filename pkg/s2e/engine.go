package s2e

import (
	"fmt"
	"io"
	"reflect"

	"github.com/yourusername/go-struct2excel/internal/export"
	"github.com/yourusername/go-struct2excel/internal/schema"
)

// Builder orchestrates the Struct-to-Excel conversion pipeline using Generics.
type Builder[T any] struct {
	outWriter io.Writer
	format    export.Format
}

// New creates a new engine for a specific struct type.
func New[T any]() *Builder[T] {
	return &Builder[T]{
		format: export.FormatXLSX,
	}
}

func (b *Builder[T]) SetOutputWriter(w io.Writer) *Builder[T] {
	b.outWriter = w
	return b
}

func (b *Builder[T]) SetFormat(ext string) *Builder[T] {
	b.format = export.Format(ext)
	return b
}

// Generate processes a slice of structs, extracts json tags as headers, and writes data.
func (b *Builder[T]) Generate(data []T) error {
	if b.outWriter == nil {
		return fmt.Errorf("output writer must be set")
	}

	if len(data) == 0 {
		return fmt.Errorf("data slice is empty")
	}

	exporter, err := export.NewExporter(b.format)
	if err != nil {
		return err
	}

	// 1. Analyze the struct tags via Reflection ONLY ONCE for performance
	dataType := reflect.TypeOf(data[0])
	cachedSchema := schema.Parse(dataType)

	// 2. Write Headers based on `json:"tag_name"`
	if err := exporter.Init(cachedSchema.Headers); err != nil {
		return fmt.Errorf("failed to initialize headers: %w", err)
	}

	// 3. Process Rows (Extremely fast due to cached schema indices)
	for _, item := range data {
		val := reflect.ValueOf(item)
		rowValues := cachedSchema.ExtractRowValues(val)

		if err := exporter.WriteRow(rowValues); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	// 4. Stream output to the destination
	return exporter.ExportTo(b.outWriter)
}
