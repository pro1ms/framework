package main

type Dependency struct {
	Name        string
	Func        string
	PackagePath string
	PackageName string
	Props       []Property
}

type Property struct {
	Name        string
	Alias       string
	IsExported  bool
	ImportPath  string
	ImportAlias string
	Inject      string
}
