package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// DBToCsvParser writes a database client's tabular output as CSV.
type DBToCsvParser interface {
	Execute(content string) error
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	content, err := readInput()
	if err != nil {
		return fmt.Errorf("read standard input: %w", err)
	}

	fileName, err := getOutputFileName()
	if err != nil {
		return fmt.Errorf("choose output file: %w", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create CSV file: %w", err)
	}
	defer file.Close()

	parser, err := NewParserForOutput(file, content)
	if err != nil {
		return err
	}
	if err := parser.Execute(content); err != nil {
		return fmt.Errorf("convert database output: %w", err)
	}

	fmt.Println(fileName)
	return nil
}

// getOutputFileName returns a unique CSV path in the user's Documents directory.
func getOutputFileName() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Documents", uuid.NewString()+".csv"), nil
}

// readInput reads all text piped to the program from standard input.
func readInput() (string, error) {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// NewParserForOutput returns the parser selected for the supplied database output.
func NewParserForOutput(w io.Writer, content string) (DBToCsvParser, error) {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	switch {
	case isPostgresFormat(lines):
		return NewPostgreParser(w), nil
	case isMySQLFormat(lines):
		return nil, nil
	default:
		return nil, errors.New("unsupported database output format")
	}
}

func isPostgresFormat(lines []string) bool {
	if len(lines) == 0 {
		return false
	}
	rowCount := regexp.MustCompile(`^\(\d+ rows?\)$`)
	return rowCount.MatchString(strings.TrimSpace(lines[len(lines)-1]))
}

func isMySQLFormat(lines []string) bool {
	return false
}
