package schema

import (
	"reflect"
	"strings"
	"sync"
)

// globalSchemaCache ensures O(1) reflection overhead for high-traffic environments.
var globalSchemaCache sync.Map

type CachedSchema struct {
	Headers    []string
	FieldPaths [][]int // Stores the exact path to handle nested structs recursively
}

// Parse extracts JSON tags and caches the schema globally.
func Parse(t reflect.Type) *CachedSchema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 1. Check Global Cache First
	if cached, ok := globalSchemaCache.Load(t); ok {
		return cached.(*CachedSchema)
	}

	var headers []string
	var paths [][]int

	// Recursive function to flatten nested structs
	var extract func(typ reflect.Type, currentPath []int)
	extract = func(typ reflect.Type, currentPath []int) {
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			jsonTag := field.Tag.Get("json")

			// Skip unexported fields or explicitly ignored json tags
			if !field.IsExported() || jsonTag == "-" {
				continue
			}

			// If the field is a nested struct (and not a time.Time), traverse into it
			if field.Type.Kind() == reflect.Struct && field.Type.String() != "time.Time" {
				newPath := append(append([]int(nil), currentPath...), i)
				extract(field.Type, newPath)
				continue
			}

			// Extract header from JSON tag (natively supports both English and Persian UTF-8)
			headerName := field.Name
			if jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				if parts[0] != "" {
					headerName = parts[0] // e.g., "firstName" or "نام کوچک"
				}
			}

			headers = append(headers, headerName)
			paths = append(paths, append(append([]int(nil), currentPath...), i))
		}
	}

	// Start extracting from the root struct
	extract(t, nil)

	schema := &CachedSchema{
		Headers:    headers,
		FieldPaths: paths,
	}

	// Store in Global Cache safely
	globalSchemaCache.Store(t, schema)
	return schema
}

// ExtractRowValues uses the cached paths to rapidly extract values, even from deep nested structs.
func (c *CachedSchema) ExtractRowValues(v reflect.Value) []interface{} {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	row := make([]interface{}, len(c.FieldPaths))
	for i, path := range c.FieldPaths {
		row[i] = v.FieldByIndex(path).Interface()
	}
	return row
}
