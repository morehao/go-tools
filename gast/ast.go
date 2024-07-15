package gast

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
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

// addMethodToInterfaceInFile reads the specified Go file, finds the specified interface,
// and adds the specified method to it if it's not already there.
func addMethodToInterfaceInFile(filePath, interfaceName, methodName string) error {
	// Read the source file
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Parse the source code into an AST
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return err
	}

	// Track whether the method is already in the interface
	methodExists := false

	// Find the specified interface and check if the method exists
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == interfaceName {
				iface, ok := x.Type.(*ast.InterfaceType)
				if !ok {
					// This is not an interface, skip it
					return true
				}
				for _, method := range iface.Methods.List {
					if len(method.Names) > 0 && method.Names[0].Name == methodName {
						methodExists = true
						return false // Stop inspecting, we found the method
					}
				}
			}
		}
		return true
	})

	// If the method doesn't exist, add it to the interface
	if !methodExists {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				if x.Name.Name == interfaceName {
					iface, ok := x.Type.(*ast.InterfaceType)
					if !ok {
						return true
					}
					iface.Methods.List = append(iface.Methods.List, &ast.Field{
						Names: []*ast.Ident{ast.NewIdent(methodName)},
						Type: &ast.FuncType{
							Params:  &ast.FieldList{},
							Results: &ast.FieldList{List: []*ast.Field{{Type: ast.NewIdent("string")}}},
						},
					})
				}
			}
			return true
		})
	}

	// If we made changes, write the file back
	if !methodExists {
		// Generate the new code from the modified AST
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, file); err != nil {
			return err
		}

		// Write the new code back to the file
		if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
			return err
		}
	}

	return nil
}
