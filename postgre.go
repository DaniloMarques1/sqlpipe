package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

type postgreParser struct {
	w *csv.Writer
}

func NewPostgreParser(file io.Writer) *postgreParser {
	return &postgreParser{w: csv.NewWriter(file)}
}

func (p *postgreParser) Execute(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("empty psql output")
	}

	lines := strings.Split(content, "\n")
	if len(lines) < 2 {
		return errors.New("psql output is missing its header separator")
	}

	err := p.parseHeader(lines[0])
	if err != nil {
		return err
	}

	if err := p.parseRows(lines[2:]); err != nil {
		return err
	}
	p.w.Flush()
	return p.w.Error()
}

func (p *postgreParser) parseHeader(header string) error {
	headerContent := strings.Split(header, "|")
	headerContentTrimmed := make([]string, 0, len(headerContent))
	for _, h := range headerContent {
		headerContentTrimmed = append(headerContentTrimmed, strings.TrimSpace(h))
	}

	if err := p.w.Write(headerContentTrimmed); err != nil {
		return err
	}

	return nil
}

func (p *postgreParser) parseRows(lines []string) error {
	for _, line := range lines {
		if !strings.Contains(line, "|") {
			continue
		}

		row := strings.Split(line, "|")
		rowTrimmed := make([]string, 0, len(row))
		for _, r := range row {
			rowTrimmed = append(rowTrimmed, strings.TrimSpace(r))
		}

		if err := p.w.Write(rowTrimmed); err != nil {
			return err
		}
	}

	return nil
}
