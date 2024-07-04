package gast

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func ParseFile(file string) (interface{}, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	serviceVars := make(map[string]string) // 用于存储service实例化变量

	// 提取文件中所有service变量的实例化
	ast.Inspect(node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						if valueSpec.Values != nil {
							if callExpr, ok := valueSpec.Values[0].(*ast.SelectorExpr); ok {
								serviceVars[name.Name] = callExpr.Sel.Name
							}
						}
					}
				}
			}
		}
		return true
	})
	return nil, err
}
