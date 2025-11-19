package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadURLsFromFile reads URLs from a text file (one per line, ignoring empty lines and comments)
func LoadURLsFromFile(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("file '%s' not found: %w", filepath, err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file '%s': %w", filepath, err)
	}

	return urls, nil
}
