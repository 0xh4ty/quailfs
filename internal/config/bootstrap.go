package config

import (
	"bufio"
	"os"
	"strings"
)

func ReadBootstrapConfig(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bootstrapNodes []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		bootstrapNodes = append(bootstrapNodes, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return bootstrapNodes, nil
}
