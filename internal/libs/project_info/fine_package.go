package project_info

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func FindPackage(packageName string, dir string) (string, error) {
	info, err := GetProjectInfo(dir)
	if err != nil {
		return "", err
	}

	var foundDir string

	// Рекурсивно обходим проект от корня (RootPath из go.mod)
	err = filepath.WalkDir(info.RootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем скрытые папки (типа .git) и vendor
		if d.IsDir() && (strings.HasPrefix(d.Name(), ".") || d.Name() == "vendor" || d.Name() == "node_modules") {
			return filepath.SkipDir
		}

		// Если нашли папку с нужным именем (например, "di")
		if d.IsDir() && d.Name() == packageName {
			foundDir = path
			return filepath.SkipAll // Нашли, выходим из поиска
		}

		return nil
	})

	if err != nil && err != filepath.SkipAll {
		return "", err
	}

	if foundDir == "" {
		return "", fmt.Errorf("package '%s' not found in project %s", packageName, info.RootPath)
	}

	// Вычисляем относительный путь от корня проекта (там где go.mod) до найденной папки
	relPath, err := filepath.Rel(info.RootPath, foundDir)
	if err != nil {
		return "", err
	}

	// Склеиваем ModuleName и относительный путь через прямой слеш (для всех ОС)
	importPath := filepath.ToSlash(filepath.Join(info.ModuleName, relPath))

	return importPath, nil
}
