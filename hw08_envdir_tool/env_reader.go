package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ReadDir: %w", err)
	}

	envResult := make(Environment)

	for _, fileInfo := range files {
		fileName := fileInfo.Name()
		isValidFileName := !strings.Contains(fileName, "=")
		if !isValidFileName {
			continue
		}

		filePath := filepath.Join(dir, fileName)
		lineStr, err := readFirstLineFile(filePath)
		if err != nil {
			continue
		}

		envResult[fileName] = getEnvValue(lineStr)
	}

	return envResult, nil
}

func readFirstLineFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("readFirstLineFile: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	firstLine, err := reader.ReadString('\n')
	if err != nil && !errors.Is(io.EOF, err) {
		return "", fmt.Errorf("readFirstLineFile: %w", err)
	}

	return firstLine, nil
}

func getEnvValue(envStr string) EnvValue {
	if len(envStr) == 0 {
		return EnvValue{"", true}
	}

	trimedStr := strings.TrimRight(envStr, " \t\n")
	res := strings.ReplaceAll(trimedStr, "\x00", "\n")

	return EnvValue{res, false}
}
