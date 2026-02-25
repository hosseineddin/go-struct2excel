# Go Struct to Excel/CSV (s2e) 🚀

[![Go CI](https://github.com/hosseineddin/go-struct2excel/actions/workflows/ci.yml/badge.svg)](https://github.com/hosseineddin/go-struct2excel/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hosseineddin/go-struct2excel)](https://goreportcard.com/report/github.com/hosseineddin/go-struct2excel)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hosseineddin/go-struct2excel)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A blazing-fast, enterprise-grade, zero-copy Go library that converts slices of deeply nested Structs directly into `XLSX` (Excel) or `CSV` files. Designed natively for the **Gin Web Framework**, it streams data directly to the client's network socket, ensuring flat RAM usage even when exporting millions of rows.

## 🧠 Enterprise Architecture Highlights

* **Native Gin Integration:** Binds directly to `*gin.Context` to handle HTTP headers and streaming automatically.
* **O(1) Reflection Overhead:** Utilizes a global thread-safe cache (`sync.Map`). Struct schemas are parsed only once per server lifecycle.
* **Massive Scale & Auto-Pagination:** Successfully tested with 2,000,000+ rows. Automatically splits data into `Sheet1`, `Sheet2`, etc., every 1,000,000 rows.
* **Deep Nested Structs:** Recursively flattens complex nested models seamlessly.
* **Native UTF-8 JSON Tags:** Uses standard `json` tags for column headers, perfectly supporting UTF-8 (e.g., Persian/Arabic characters) without custom tag syntax.
* **Secure by Default:** Respects `json:"-"` and unexported fields to prevent data leaks (e.g., passwords or internal IDs).

## 📦 Installation

```bash
go get [github.com/hosseineddin/go-struct2excel](https://github.com/hosseineddin/go-struct2excel)