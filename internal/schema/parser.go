package schema

import (
	"reflect"
	"strings"
)

// CachedSchema holds the structure of the struct to avoid reflecting every row
type CachedSchema struct {
	Headers    []string
	FieldIndex []int
}

// Parse extracts the json tags from a struct to use as Excel headers.
// It skips fields with json:"-" or fields without a json tag (optional based on your need).
func Parse(t reflect.Type) *CachedSchema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var headers []string
	var indices []int

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")

		if tag == "-" {
			continue // Skip ignored fields
		}

		headerName := field.Name // Default to struct field name
		if tag != "" {
			// Handle tags like `json:"first_name,omitempty"`
			parts := strings.Split(tag, ",")
			if parts[0] != "" {
				headerName = parts[0]
			}
		}

		headers = append(headers, headerName)
		indices = append(indices, i)
	}

	return &CachedSchema{
		Headers:    headers,
		FieldIndex: indices,
	}
}

// ExtractRowValues uses the cached schema to rapidly extract values from a struct
func (c *CachedSchema) ExtractRowValues(v reflect.Value) []interface{} {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	row := make([]interface{}, len(c.FieldIndex))
	for i, fieldIdx := range c.FieldIndex {
		row[i] = v.Field(fieldIdx).Interface()
	}
	return row
}
