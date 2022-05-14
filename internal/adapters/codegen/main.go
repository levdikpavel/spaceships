package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
)

var (
	inputFile  string
	outputFile string
	interfaces string
)

func init() {
	flag.StringVar(&inputFile, "f", "", "input file")
	flag.StringVar(&outputFile, "o", "", "output file")
	flag.StringVar(&interfaces, "i", "", "interfaces to generate adapters")
}

func main() {
	flag.Parse()
	splittedInterfaces := strings.Split(interfaces, ",")

	processor := newAdapterProcessor()
	data := processor.parseFile(inputFile, splittedInterfaces)
	dumpOutput(data, outputFile)
}

func dumpOutput(data templateStruct, outputFile string) {
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("test").Parse(outputTemplate)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		log.Fatal(err)
	}
}

type adaptersProcessor struct {
	imports map[string]string
	typesImports map[string]struct{}
}

func newAdapterProcessor() *adaptersProcessor {
	return &adaptersProcessor{
		imports: map[string]string{
			"ioc":  `"modules/internal/ioc"`,
			"core": `"modules/internal/core"`,
		},
		typesImports: map[string]struct{}{
			"ioc":  {},
			"core": {},
		},
	}
}

func (a *adaptersProcessor) parseFile(path string, interfaces []string) templateStruct {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	interfacesMap := make(map[string]struct{})
	for _, interfaceStr := range interfaces {
		interfacesMap[interfaceStr] = struct{}{}
	}

	parsedInterfaces := a.parseNode(node, interfacesMap)

	var imports []string
	for importAlias, _ := range a.typesImports {
		importValue, ok := a.imports[importAlias]
		if !ok {
			fmt.Printf("missing import '%s'\n", importAlias)
			continue
		}
		imports = append(imports, importValue)
	}

	sort.Strings(imports)

	return templateStruct{
		Imports:    imports,
		Interfaces: parsedInterfaces,
	}
}

func (a *adaptersProcessor) parseNode(
	node *ast.File, interfaces map[string]struct{},
) (outInterfaces []interfaceDecl) {
	for _, f := range node.Decls {
		genD, ok := f.(*ast.GenDecl)
		if !ok {
			fmt.Printf("SKIP %T is not *ast.GenDecl\n", f)
			continue
		}

		for _, spec := range genD.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if ok {
				a.parseImport(importSpec)
				continue
			}

			currentType, ok := spec.(*ast.TypeSpec)
			if !ok {
				fmt.Printf("SKIP %T is not ast.TypeSpec\n", spec)
				continue
			}

			currentInterface, ok := currentType.Type.(*ast.InterfaceType)
			if !ok {
				fmt.Printf("SKIP %T is not ast.InterfaceType\n", currentInterface)
				continue
			}

			interfaceName := currentType.Name.Name
			if _, ok := interfaces[interfaceName]; ok {
				outInterface := a.parseInterface(interfaceName, currentInterface.Methods.List)
				outInterfaces = append(outInterfaces, outInterface)
			}
		}
	}
	return
}

func (a *adaptersProcessor) parseImport(importSpec *ast.ImportSpec) {
	value := importSpec.Path.Value
	trimmed := strings.Trim(value, `"`)
	splitted := strings.Split(trimmed, "/")
	alias := splitted[len(splitted) - 1]
	a.imports[alias] = value
}

func (a *adaptersProcessor) parseInterface(name string, fields []*ast.Field) interfaceDecl {
	var methods []methodDecl
	for _, field := range fields {
		method := a.parseMethod(field)
		methods = append(methods, method)
	}

	return interfaceDecl{
		Name:    name,
		Methods: methods,
	}
}

func (a *adaptersProcessor) parseMethod(field *ast.Field) methodDecl {
	name := field.Names[0].Name
	var property string
	var set bool
	switch {
	case strings.HasPrefix(name, "Get"):
		property = strings.TrimPrefix(name, "Get")
	case strings.HasPrefix(name, "Set"):
		property = strings.TrimPrefix(name, "Set")
		set = true
	default:
		panic("unsupported method name " + name)
	}
	var params []ParamDecl
	var results []ParamDecl
	for _, funcParam := range field.Type.(*ast.FuncType).Params.List {
		params = append(params, a.parseParam(funcParam))
	}
	for _, funcResult := range field.Type.(*ast.FuncType).Results.List {
		results = append(results, a.parseParam(funcResult))
	}

	return methodDecl{
		Name:     name,
		Property: property,
		Set:      set,
		Params:   params,
		Results:  results,
	}
}

func (a *adaptersProcessor) parseParam(field *ast.Field) ParamDecl {
	var packageName string
	var typeName string

	switch t := field.Type.(type) {
	case *ast.SelectorExpr:
		packageName = t.X.(*ast.Ident).Name
		typeName = t.Sel.Name
	case *ast.Ident:
		typeName = t.Name
	default:
		panic(t)
	}

	paramType := typeName
	if packageName != "" {
		paramType = packageName + "." + typeName
	}

	a.typesImports[packageName] = struct{}{}

	return ParamDecl{
		Name: "",
		Type: paramType,
	}
}

type templateStruct struct {
	Imports    []string
	Interfaces []interfaceDecl
}

type interfaceDecl struct {
	Name    string
	Methods []methodDecl
}

type methodDecl struct {
	Name     string
	Property string
	Set      bool
	Params   []ParamDecl
	Results  []ParamDecl
}

type ParamDecl struct {
	Name string
	Type string
}

const (
	outputTemplate = `package adapters

import ({{range .Imports}}
	{{.}}{{end}}
)

{{range $interface := .Interfaces}}
type {{.Name}}Adapter struct{
   obj core.UObject
}

func New{{.Name}}Adapter(obj core.UObject) *{{.Name}}Adapter {
	return &{{.Name}}Adapter{
		obj: obj,
	}
}
{{range $method := .Methods}}
func (a *{{$interface.Name}}Adapter) {{.Name}}({{range .Params}}newValue {{.Type}}{{end}}) ({{range $i, $v := .Results}}{{if $i}}, {{end}}{{.Type}}{{end}}) {
{{- if $method.Set}}
	return ioc.Resolve("Operations.{{$interface.Name}}:{{$method.Property}}.set", a.obj, newValue).(core.Command).Execute()
{{- else}}
	return ioc.Resolve("Operations.{{$interface.Name}}:{{$method.Property}}.get", a.obj).({{(index .Results 0).Type}}), nil
{{- end}}
}
{{end}}
{{- end}}

func init() {
	_ = ioc.Resolve("IoC.Register", "Adapter",
		func(params ...interface{}) interface{} {
			adapterType := params[0].(string)
			switch adapterType {
{{- range $interface := .Interfaces}}
			case "core.{{.Name}}":
				return New{{.Name}}Adapter(params[1].(core.UObject))
{{- end}}
			default:
				panic("unknown adapter type" + adapterType)
			}
	}).(core.Command).Execute()
}
`
)
