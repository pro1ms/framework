package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func generate(outDir string, deps []Dependency) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	// 2. Создаем базовый файл container.go (если он не существует)
	baseFile := filepath.Join(outDir, "container.gen.go")
	if _, err := os.Stat(baseFile); os.IsNotExist(err) {
		_ = os.WriteFile(baseFile, []byte(baseCode), 0644)
	}

	// 3. Генерируем файлы для каждой зависимости
	tmpl := template.Must(template.New("service").Parse(serviceTemplate))

	for _, d := range deps {
		err := generateFile(outDir, d, tmpl)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateFile(outDir string, d Dependency, tmpl *template.Template) error {
	fileName := fmt.Sprintf("%s.gen.go", strings.ToLower(d.Name))
	f, err := os.Create(filepath.Join(outDir, fileName))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	// Подготовка данных: вычленяем тип возвращаемого значения (упрощенно)
	// Для NewHelloHandler -> HelloHandler
	returnType := strings.TrimPrefix(d.Func, "New")

	// Собираем уникальные импорты для аргументов (Props)
	type Import struct{ Alias, Path string }
	var extraImports []Import
	// Тут можно добавить логику сбора импортов из Property, если они из других пакетов

	data := struct {
		Dependency
		Type    string
		Imports []Import
	}{
		Dependency: d,
		Type:       returnType,
		Imports:    extraImports,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}
