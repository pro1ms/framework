package http

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pro1ms/framework/internal/libs/names"
)

func (s *Scanner) write(destDir string) error {
	err := s.writeHandlers(destDir)
	if err != nil {
		return err
	}
	err = s.writeRouter(destDir)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scanner) writeGetPackagePath(dir string, pkg Package) string {
	return filepath.Join(dir, pkg.Name)
}

func (s *Scanner) writeGetMethodFile(dir string, pkg Package, method Method) string {
	pkgPath := s.writeGetPackagePath(dir, pkg)
	return filepath.Join(pkgPath, names.ToSnakeCase(method.Name)+".go")
}

func (s *Scanner) writeHandlers(destDir string) error {
	for _, pkg := range s.Packages {
		modulePath := s.writeGetPackagePath(destDir, pkg)
		err := os.MkdirAll(modulePath, 0755)
		if err != nil {
			return err
		}

		for _, method := range pkg.Methods {
			fullPath := s.writeGetMethodFile(destDir, pkg, method)

			// 3. Проверка: если файл уже есть, не трогаем его!
			if _, err := os.Stat(fullPath); err == nil {
				continue
			}

			// 4. Подготовка данных для шаблона
			data := struct {
				PackageName string
				ImportPath  string
				HandlerName string
				Methods     []Method
			}{
				PackageName: pkg.Name,
				ImportPath:  pkg.Path,
				HandlerName: method.Name, // Например: GetHelloHandler
				Methods:     []Method{method},
			}

			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				err := s.writeTemplate(fullPath, data, handlerTemplate)
				if err != nil {
					return fmt.Errorf("failed to write handler %s: %w", fullPath, err)
				}
				fmt.Printf("Gen HTTP: write handler - %s\n", fullPath)
			} else {
				fmt.Printf("Gen HTTP: skip existing handler %s\n", method.Name)
			}
		}
	}
	return nil
}

func (s *Scanner) writeTemplate(fullPath string, data any, tpl string) error {
	funcMap := template.FuncMap{
		"printArgs": func(m Method) string {
			var args []string
			for _, a := range m.Args {
				args = append(args, fmt.Sprintf("%s %s", a.Name, a.Type))
			}
			return strings.Join(args, ", ")
		},
		"printResults": func(m Method) string {
			if len(m.Results) > 1 {
				return "(" + strings.Join(m.Results, ", ") + ")"
			}
			return strings.Join(m.Results, ", ")
		},
	}

	// Теперь используем аргумент tpl вместо константы
	t, err := template.New("gen").Funcs(funcMap).Parse(tpl)
	if err != nil {
		return fmt.Errorf("failed parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed execute template: %w", err)
	}

	// Форматируем Go-код
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return os.WriteFile(fullPath, buf.Bytes(), 0644)
	}

	return os.WriteFile(fullPath, formatted, 0644)
}
