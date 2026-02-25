package export

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

const maxRowsPerSheet = 1000000

type xlsxExporter struct {
	file         *excelize.File
	streamWriter *excelize.StreamWriter
	sheetIndex   int
	currentRow   int
	headers      []string
}

func newXLSXExporter() *xlsxExporter {
	return &xlsxExporter{file: excelize.NewFile(), sheetIndex: 1, currentRow: 1}
}

func (x *xlsxExporter) Init(headers []string) error {
	x.headers = headers
	sheetName := fmt.Sprintf("Sheet%d", x.sheetIndex)
	if x.sheetIndex > 1 {
		x.file.NewSheet(sheetName)
	} else {
		x.file.SetSheetName("Sheet1", sheetName)
	}

	sw, err := x.file.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}
	x.streamWriter = sw

	headerRow := make([]interface{}, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	cell, _ := excelize.CoordinatesToCellName(1, 1)
	x.streamWriter.SetRow(cell, headerRow)
	x.currentRow = 2
	return nil
}

func (x *xlsxExporter) WriteRow(row []interface{}) error {
	if x.currentRow > maxRowsPerSheet {
		x.streamWriter.Flush()
		x.sheetIndex++
		x.currentRow = 1
		x.Init(x.headers)
	}
	cell, _ := excelize.CoordinatesToCellName(1, x.currentRow)
	x.streamWriter.SetRow(cell, row)
	x.currentRow++
	return nil
}

func (x *xlsxExporter) ExportTo(w io.Writer) error {
	x.streamWriter.Flush()
	return x.file.Write(w)
}
