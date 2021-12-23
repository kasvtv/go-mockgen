package generation

import (
	"go/ast"

	"github.com/dave/jennifer/jen"
)

func compose(stmt *jen.Statement, tail ...jen.Code) *jen.Statement {
	head := *stmt
	for _, value := range tail {
		head = append(head, value)
	}

	return &head
}

func addComment(code *jen.Statement, level int, commentText string) *jen.Statement {
	if commentText == "" {
		return code
	}

	comment := generateComment(level, commentText)
	return compose(comment, code)
}

func addTypes(code *jen.Statement, fields []*ast.Field, includeTypes bool) *jen.Statement {
	if len(fields) == 0 {
		return code
	}

	types := make([]jen.Code, 0, len(fields))
	for _, field := range fields {
		for _, name := range field.Names {
			if includeTypes {
				// TODO - use actual constraint
				types = append(types, jen.Id(name.Name).Any())
			} else {
				types = append(types, jen.Id(name.Name))
			}
		}
	}

	return compose(code, jen.Types(types...))
}

func selfAppend(sliceRef *jen.Statement, value jen.Code) jen.Code {
	return compose(sliceRef, jen.Op("=").Id("append").Call(sliceRef, value))
}
