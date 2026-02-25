# Go Struct to Excel/CSV (s2e) 

[![Go CI](https://github.com/hosseineddin/go-struct2excel/actions/workflows/ci.yml/badge.svg)](https://github.com/hosseineddin/go-struct2excel/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hosseineddin/go-struct2excel)](https://goreportcard.com/report/github.com/hosseineddin/go-struct2excel)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hosseineddin/go-struct2excel)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A blazing-fast, type-safe, and reflection-optimized Go library that converts slices of Go Structs directly into `XLSX` (Excel) or `CSV` files. It dynamically reads standard `json:"tag"` annotations to generate Excel columns, making it an instant plug-and-play solution for existing backend models (e.g., GORM models or API responses).

## Architecture Highlights
Unlike naive approaches that use Reflection on every single row (causing massive CPU spikes), `go-struct2excel` uses **Schema Caching**:
1. **Generics `[T any]`**: Ensures type safety at compile time.
2. **O(1) Reflection Overhead**: The engine parses the struct's `json` tags *only once* for the entire dataset and caches the field indices. Processing 1 million records incurs almost zero reflection penalty.
3. **Zero-Copy Streaming**: Writes the generated Excel/CSV file directly to an `io.Writer` (like an HTTP Response), keeping server RAM usage flat.
4. **Smart Tag Parsing**: Respects `json:"-"` to ignore sensitive fields (like passwords) and extracts the correct column names from tags like `json:"first_name,omitempty"`.

## Installation
```bash
go get [github.com/hosseineddin/go-struct2excel](https://github.com/hosseineddin/go-struct2excel)