package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

type DBToCsv interface {
	Execute(content string) error
}

type postgreDBToCsv struct {
	w *csv.Writer
}

func NewPostgreDBToCsv(file io.Writer) *postgreDBToCsv {
	return &postgreDBToCsv{w: csv.NewWriter(file)}
}

func (p *postgreDBToCsv) Execute(content string) error {
	defer p.w.Flush()

	lines := strings.Split(content, "\n")
	header, err := p.parseHeader(lines[0])
	if err != nil {
		return err
	}

	if err := p.w.Write(header); err != nil {
		return err
	}

	if err := p.parseRows(lines[2:]); err != nil {
		return err
	}

	return nil
}

func (p *postgreDBToCsv) parseHeader(header string) ([]string, error) {
	headerContent := strings.Split(header, "|")
	headerContentTrimmed := make([]string, 0, len(headerContent))
	for _, h := range headerContent {
		headerContentTrimmed = append(headerContentTrimmed, strings.TrimSpace(h))
	}
	return headerContentTrimmed, nil
}

func (p *postgreDBToCsv) parseRows(lines []string) error {
	for _, line := range lines {
		if !strings.Contains(line, "|") {
			continue
		}

		row := strings.Split(line, "|")
		rowTrimmed := make([]string, 0, len(row))
		for _, r := range row {
			trimmed := strings.TrimSpace(r)
			if len(trimmed) > 0 {
				rowTrimmed = append(rowTrimmed, trimmed)
			}
		}

		if err := p.w.Write(rowTrimmed); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	content, err := readInput()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := NewPostgreDBToCsv(f)
	p.Execute(content)

}

func readInput() (string, error) {
	b, err := os.ReadFile("test.txt") // TODO: we just receive from stdin
	if err != nil {
		return "", err
	}

	return string(b), nil
}
