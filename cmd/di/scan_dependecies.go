package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
)

func scanDependencies(root string) ([]Dependency, error) {
	var dependencies []Dependency

	absRoot, _ := filepath.Abs(root)

	fmt.Println("DI: start scan")

	fileSet := token.NewFileSet()
	err := filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// 1. Вычисляем путь к папке, где лежит текущий файл
		dir := filepath.Dir(path)
		// 2. Вычисляем относительный путь от корня проекта
		relDir, _ := filepath.Rel(absRoot, dir)
		// 3. Формируем полный Import Path для Go
		packagePath := filepath.Join("github.com/gor/example/internal", relDir)
		packagePath = filepath.ToSlash(packagePath)

		file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("di: failed parse file %s - %w", path, err)
		}

		deps, err := inspectFile(packagePath, fileSet, file, path)
		if err != nil {
			return fmt.Errorf("di: failed inspect file %s - %w", path, err)
		}
		dependencies = append(dependencies, deps...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	fmt.Println("DI: finished scan")
	return dependencies, nil
}
