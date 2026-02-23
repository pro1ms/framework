package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gor/framework/internal/libs/project_info"
)

type Scanner struct {
	ProjectInfo project_info.ProjectInfo
	Packages    []Package
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Run(in string, out string) error {
	projectInfo, err := project_info.GetProjectInfo(in)
	if err != nil {
		fmt.Println("Gen HTTP: failed to get project info:")
		fmt.Println(err)
		return err
	}
	s.ProjectInfo = *projectInfo

	fmt.Println("Gen HTTP: scan start")
	err = s.scan(in)
	if err != nil {
		fmt.Println("Gen HTTP: scan failed:")
		fmt.Println(err)
		return err
	}
	fmt.Println("Gen HTTP: scan finished")

	fmt.Println("Gen HTTP: write start")
	err = s.write(out)
	if err != nil {
		fmt.Println("Gen HTTP: write failed:")
		fmt.Println(err)
		return err
	}
	fmt.Println("Gen HTTP: write finished")

	fmt.Println("Gen HTTP: compilation start")
	err = s.checkCompilation(out)
	if err != nil {
		return err
	}
	fmt.Println("Gen HTTP: compilation finished")

	return nil
}

func (s *Scanner) scan(dir string) error {
	root, _ := filepath.Abs(dir)
	fileSet := token.NewFileSet()
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("gen http: failed parse file %s - %w", path, err)
		}

		info, err := s.resolveImportPath(path)
		if err != nil {
			return fmt.Errorf("failed resolve package: %w", err)
		}

		methods, err := s.scanFile(file, *info)
		if err != nil {
			return fmt.Errorf("gen http: failed scan file %s - %w", path, err)
		}

		if len(methods) > 0 {

			s.Packages = append(s.Packages, Package{
				Name:    info.Name,
				Path:    info.Path,
				Methods: methods,
			})
		}

		return err
	})
}

func (s *Scanner) scanFile(file *ast.File, pkg PackageInfo) ([]Method, error) {
	for _, decl := range file.Decls {
		// Ищем общие объявления (type, const, var)
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// Ищем интерфейс с конкретным именем
			if typeSpec.Name.Name == "StrictServerInterface" {
				interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
				if !ok {
					return nil, fmt.Errorf("StrictServerInterface is not an interface")
				}

				methods := make([]Method, 0)

				// Перебираем методы интерфейса
				for _, method := range interfaceType.Methods.List {
					// Имя метода (может быть несколько, если в одной строке, но обычно одно)
					for _, methodName := range method.Names {
						// Тут можно вытащить сигнатуру (аргументы и возвращаемые значения)
						if funcType, ok := method.Type.(*ast.FuncType); ok {
							methods = append(methods, s.analyzeMethod(methodName.Name, pkg, funcType))
						}
					}
				}
				return methods, nil
			}
		}
	}

	return nil, nil
}

func (s *Scanner) analyzeMethod(name string, pkg PackageInfo, f *ast.FuncType) Method {
	fmt.Printf("Gen HTTP: found method - %s\n", name)

	method := Method{
		Name: name,
	}

	// 1. Парсим входящие аргументы (Params)
	if f.Params != nil {
		for _, field := range f.Params.List {
			typeStr := s.nodeToString(field.Type, pkg.Name)
			// У аргумента может быть несколько имен (a, b int)
			if len(field.Names) > 0 {
				for _, n := range field.Names {
					method.Args = append(method.Args, MethodArg{Name: n.Name, Type: typeStr})
				}
			} else {
				// Если имя не указано (анонимный аргумент в интерфейсе)
				method.Args = append(method.Args, MethodArg{Name: "", Type: typeStr})
			}
		}
	}

	// 2. Парсим возвращаемые значения (Results)
	if f.Results != nil {
		for _, field := range f.Results.List {
			method.Results = append(method.Results, s.nodeToString(field.Type, pkg.Name))
		}
	}

	return method
}

// Вспомогательный метод для превращения AST-узла в строку (название типа)
func (s *Scanner) nodeToString(node ast.Node, pkgName string) string {
	var buf bytes.Buffer
	// token.NewFileSet() здесь ок, так как нам нужна просто строка
	if err := format.Node(&buf, token.NewFileSet(), node); err != nil {
		return ""
	}

	typeStr := buf.String()

	// Список типов, которые НЕ нужно трогать
	builtInTypes := map[string]bool{
		"error":           true,
		"string":          true,
		"int":             true,
		"bool":            true,
		"context.Context": true, // context уже с точкой, его не трогаем
	}

	// Если тип не содержит точку и не является встроенным — добавляем префикс
	if !strings.Contains(typeStr, ".") && !builtInTypes[typeStr] {
		return pkgName + "." + typeStr
	}

	return typeStr
}

func (s *Scanner) resolveImportPath(filePath string) (*PackageInfo, error) {
	info, err := project_info.GetProjectInfo(filePath)
	if err != nil {
		return nil, err
	}

	absFile, _ := filepath.Abs(filePath)
	dir := filepath.Dir(absFile)

	// Получаем относительный путь от корня проекта до папки с файлом
	relPath, err := filepath.Rel(info.RootPath, dir)
	if err != nil {
		return nil, err
	}

	// Склеиваем и приводим к прямому слешу (важно для Windows)
	fullPath := filepath.ToSlash(filepath.Join(info.ModuleName, relPath))

	// В твоем случае это будет ".../internal/generated/http/api"
	// Последний элемент пути — это имя пакета
	pkgName := filepath.Base(fullPath)

	return &PackageInfo{
		Name: pkgName,
		Path: fullPath,
	}, nil
}

func (s *Scanner) printJSON() {
	// MarshalIndent делает "красивый" вывод с отступами
	data, err := json.MarshalIndent(s.Packages, "", "  ")
	if err != nil {
		fmt.Printf("failed marshal json: %v\n", err)
		return
	}

	fmt.Println(string(data))
}
