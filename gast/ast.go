package gast

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func ParseFile(file string) (interface{}, error) {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
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

// FindFunctionInFile 在指定的文件中查找函数，注意，这里只查找非方法函数
func FindFunctionInFile(file *ast.File, functionName string) (*ast.FuncDecl, bool) {
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Recv == nil {
				if funcDecl.Name.Name == functionName {
					return funcDecl, true
				}
			}
		}
	}
	return nil, false
}

// FindMethodInFile 在指定的文件中查找方法, 注意，这里只查找指定接收者类型的方法
func FindMethodInFile(file *ast.File, receiverTypeName, methodName string) (*ast.FuncDecl, bool) {
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil {
				for _, recv := range funcDecl.Recv.List {
					if recvType, ok := recv.Type.(*ast.StarExpr); ok {
						if ident, ok := recvType.X.(*ast.Ident); ok && ident.Name == receiverTypeName {
							if funcDecl.Name.Name == methodName {
								return funcDecl, true
							}
						}
					} else if ident, ok := recv.Type.(*ast.Ident); ok {
						if ident.Name == receiverTypeName {
							if funcDecl.Name.Name == methodName {
								return funcDecl, true
							}
						}
					}
				}
			}
		}
	}
	return nil, false
}

// AddMethodToInterfaceInFile 在指定的文件中查找指定的接口，并在其中添加指定的方法
func AddMethodToInterfaceInFile(filePath, interfaceName, receiverTypeName, methodName string) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fileSet := token.NewFileSet()
	file, parseFileErr := parser.ParseFile(fileSet, filePath, src, parser.ParseComments)
	if parseFileErr != nil {
		return parseFileErr
	}

	// 查找接收者类型的方法声明
	methodDecl, found := FindMethodInFile(file, receiverTypeName, methodName)
	if !found {
		return errors.New("method not found: " + methodName)
	}
	methodType := methodDecl.Type

	methodExists := false

	// 遍历文件中的所有类型声明，查找指定的接口
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == interfaceName {
				iface, ok := x.Type.(*ast.InterfaceType)
				if !ok {
					return true
				}
				for _, method := range iface.Methods.List {
					if len(method.Names) > 0 && method.Names[0].Name == methodName {
						methodExists = true
						return false
					}
				}
				if !methodExists {
					// 将方法添加到接口中
					iface.Methods.List = append(iface.Methods.List, &ast.Field{
						Names: []*ast.Ident{ast.NewIdent(methodName)},
						Type:  methodType,
					})
				}
			}
		}
		return true
	})

	if !methodExists {
		// 根据修改后的AST生成新代码
		var buf bytes.Buffer
		if err := format.Node(&buf, fileSet, file); err != nil {
			return err
		}

		// 重写文件
		if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
			return err
		}
	}

	return nil
}
