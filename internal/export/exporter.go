package export

import (
	"fmt"
	"io"
)

type Exporter interface {
	Init(headers []string) error
	WriteRow(row []interface{}) error
	ExportTo(w io.Writer) error
}

type Format string

const (
	FormatXLSX Format = "xlsx"
	FormatCSV  Format = "csv"
)

func NewExporter(format Format) (Exporter, error) {
	switch format {
	case FormatXLSX:
		return newXLSXExporter(), nil
	case FormatCSV:
		return newCSVExporter(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
