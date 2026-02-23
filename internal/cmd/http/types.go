package http

type MethodArg struct {
	Name string
	Type string
}

type PackageInfo struct {
	Name string
	Path string
}

type Method struct {
	Name    string
	Args    []MethodArg
	Results []string
}

type Package struct {
	Name    string
	Path    string
	Methods []Method
}
