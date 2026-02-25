package s2e

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/hosseineddin/go-struct2excel/internal/export"
	"github.com/hosseineddin/go-struct2excel/internal/schema"
)

// Builder orchestrates the Struct-to-Excel conversion pipeline using Generics.
type Builder[T any] struct {
	format   export.Format
	filename string
}

// New creates a new engine for a specific struct type.
func New[T any]() *Builder[T] {
	return &Builder[T]{
		format:   export.FormatXLSX,
		filename: "export",
	}
}

// SetFormat configures the output format (xlsx or csv).
func (b *Builder[T]) SetFormat(ext string) *Builder[T] {
	b.format = export.Format(ext)
	return b
}

// SetFilename sets the name of the downloaded file (without extension).
func (b *Builder[T]) SetFilename(name string) *Builder[T] {
	b.filename = name
	return b
}

// Stream directly binds to the Gin context, sets appropriate HTTP headers,
// and streams the generated Excel/CSV file directly to the client.
func (b *Builder[T]) Stream(c *gin.Context, data []T) error {
	if len(data) == 0 {
		return fmt.Errorf("data slice is empty")
	}

	// Set automatic HTTP download headers for Gin
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.%s", b.filename, b.format))
	c.Header("Content-Transfer-Encoding", "binary")

	if b.format == export.FormatCSV {
		c.Header("Content-Type", "text/csv; charset=utf-8")
	} else {
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	}

	exporter, err := export.NewExporter(b.format)
	if err != nil {
		return err
	}

	// 1. Analyze the struct tags (O(1) after first cache)
	dataType := reflect.TypeOf(data[0])
	cachedSchema := schema.Parse(dataType)

	// 2. Initialize Headers
	if err := exporter.Init(cachedSchema.Headers); err != nil {
		return fmt.Errorf("failed to initialize headers: %w", err)
	}

	// 3. Process Rows rapidly using FieldByIndex
	for _, item := range data {
		val := reflect.ValueOf(item)
		rowValues := cachedSchema.ExtractRowValues(val)

		if err := exporter.WriteRow(rowValues); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	// 4. Stream output directly to the Gin network writer
	return exporter.ExportTo(c.Writer)
}
