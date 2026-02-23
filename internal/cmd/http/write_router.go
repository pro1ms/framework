package http

import (
	"path/filepath"
	"sort"

	"github.com/pro1ms/framework/internal/libs/project_info"
)

type BridgeData struct {
	DIPath        string   // Путь до твоего пакета di
	ModuleImports []string // Пути до папок с хендлерами (internal/handlers/http/test01 и т.д.)
	Handlers      []HandlerInfo
}

type HandlerInfo struct {
	Package    string // Например: "test01"
	StructName string // Например: "GetHello"
}

func (s *Scanner) writeRouter(destDir string) error {
	filePath := filepath.Join(destDir, "router.go")
	diImport, err := project_info.FindPackage("di", filePath)
	if err != nil {
		return err
	}

	data := BridgeData{
		DIPath: diImport,
	}

	for _, pkg := range s.Packages {
		pkgPath := s.writeGetPackagePath(destDir, pkg)
		importPath, err := s.getImport(pkgPath)
		if err != nil {
			return err
		}
		data.ModuleImports = append(data.ModuleImports, importPath)

		for _, m := range pkg.Methods {
			data.Handlers = append(data.Handlers, HandlerInfo{
				Package:    pkg.Name,
				StructName: m.Name,
			})
		}
	}

	sort.Strings(data.ModuleImports)
	sort.Slice(data.Handlers, func(i, j int) bool {
		return data.Handlers[i].StructName < data.Handlers[j].StructName
	})

	return s.writeTemplate(filePath, data, routerTemplate)
}

func (s *Scanner) getImport(path string) (string, error) {
	absDest, _ := filepath.Abs(path)
	relPath, err := filepath.Rel(s.ProjectInfo.RootPath, absDest)
	if err != nil {
		return "", err
	}

	return filepath.ToSlash(filepath.Join(s.ProjectInfo.ModuleName, relPath)), nil
}
