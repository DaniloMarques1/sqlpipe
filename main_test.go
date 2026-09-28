package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"reflect"
	"testing"
)

func TestParseHeader(t *testing.T) {
	var output bytes.Buffer
	parser := NewPostgreParser(&output)
	err := parser.parseHeader("id | codigo_rastreio | cliente nome")
	if err != nil {
		t.Fatalf("parseHeader() error = %v", err)
	}

	parser.w.Flush()
	rows, err := csv.NewReader(&output).ReadAll()
	if err != nil {
		t.Fatalf("reading generated CSV: %v", err)
	}
	want := [][]string{{"id", "codigo_rastreio", "cliente nome"}}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("CSV records = %#v, want %#v", rows, want)
	}
}

func TestExecuteWritesCSVAndPreservesEmptyColumns(t *testing.T) {
	content := `
 id | cliente       | metadados
----+---------------+----------
  1 | Cliente, 1    |
  2 | Cliente 2     | ok
(2 rows)`

	var output bytes.Buffer
	parser, err := NewParserForOutput(&output, content)
	if err != nil {
		t.Fatalf("NewParserForOutput() error = %v", err)
	}
	if err := parser.Execute(content); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	rows, err := csv.NewReader(&output).ReadAll()
	if err != nil {
		t.Fatalf("reading generated CSV: %v", err)
	}
	want := [][]string{
		{"id", "cliente", "metadados"},
		{"1", "Cliente, 1", ""},
		{"2", "Cliente 2", "ok"},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("CSV records = %#v, want %#v", rows, want)
	}
}

func TestNewParserForOutput(t *testing.T) {
	if _, err := NewParserForOutput(io.Discard, "not psql output"); err == nil {
		t.Fatal("NewParserForOutput() accepted unsupported input")
	}
}
