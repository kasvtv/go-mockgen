package types

import (
	"go/ast"
	"go/types"
	"sort"
)

type Interface struct {
	Name       string
	ImportPath string
	TypeParams []*ast.Field
	Methods    []*Method
}

func newInterfaceFromTypeSpec(name, importPath string, typeSpec *ast.TypeSpec, underlyingType *types.Interface) *Interface {
	methodMap := make(map[string]*Method, underlyingType.NumMethods())
	for i := 0; i < underlyingType.NumMethods(); i++ {
		method := underlyingType.Method(i)
		name := method.Name()
		methodMap[name] = newMethodFromSignature(name, method.Type().(*types.Signature))
	}

	methodNames := make([]string, 0, len(methodMap))
	for k := range methodMap {
		methodNames = append(methodNames, k)
	}
	sort.Strings(methodNames)

	methods := make([]*Method, 0, len(methodNames))
	for _, name := range methodNames {
		methods = append(methods, methodMap[name])
	}

	var typeParams []*ast.Field
	if typeSpec.TypeParams != nil {
		typeParams = typeSpec.TypeParams.List
	}

	return &Interface{
		Name:       name,
		ImportPath: importPath,
		TypeParams: typeParams,
		Methods:    methods,
	}
}
