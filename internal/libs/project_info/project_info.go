package project_info

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProjectInfo struct {
	RootPath   string
	ModuleName string
}

func GetProjectInfo(startDir string) (*ProjectInfo, error) {
	curr, err := filepath.Abs(startDir)
	if err != nil {
		return nil, err
	}

	for {
		modPath := filepath.Join(curr, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			// Нашли go.mod, читаем имя модуля
			moduleName, err := readModuleName(modPath)
			if err != nil {
				return nil, err
			}
			return &ProjectInfo{
				RootPath:   curr,
				ModuleName: moduleName,
			}, nil
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break // Дошли до корня диска
		}
		curr = parent
	}

	return nil, fmt.Errorf("go.mod not found starting from %s", startDir)
}

func readModuleName(modPath string) (string, error) {
	file, err := os.Open(modPath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("module name not found in %s", modPath)
}
