package yesql

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type QueryLoader struct {
	queries map[string]string
}

func NewQueryLoader() *QueryLoader {
	return &QueryLoader{
		queries: make(map[string]string),
	}
}

func (ql *QueryLoader) LoadQueriesFromDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".sql" {
			return ql.loadQueriesFromFile(path)
		}

		return nil
	})
}

func (ql *QueryLoader) loadQueriesFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nameRegex := regexp.MustCompile(`^--\s*name:\s*(\w+)`)
	
	var currentQueryName string
	var currentQueryLines []string

	for scanner.Scan() {
		line := scanner.Text()
		
		if matches := nameRegex.FindStringSubmatch(line); len(matches) > 1 {
			if currentQueryName != "" {
				ql.queries[currentQueryName] = strings.Join(currentQueryLines, "\n")
			}
			
			currentQueryName = matches[1]
			currentQueryLines = []string{}
		} else if currentQueryName != "" && !strings.HasPrefix(strings.TrimSpace(line), "--") {
			currentQueryLines = append(currentQueryLines, line)
		}
	}

	if currentQueryName != "" {
		ql.queries[currentQueryName] = strings.Join(currentQueryLines, "\n")
	}

	return scanner.Err()
}

func (ql *QueryLoader) GetQuery(name string) (string, error) {
	query, exists := ql.queries[name]
	if !exists {
		return "", fmt.Errorf("query '%s' not found", name)
	}
	return strings.TrimSpace(query), nil
}

func (ql *QueryLoader) ListQueries() []string {
	var names []string
	for name := range ql.queries {
		names = append(names, name)
	}
	return names
}