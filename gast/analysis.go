package gast

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
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

func HasPackageKeywords(file string) (bool, error) {
	fileSet := token.NewFileSet()
	node, parseErr := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
	if parseErr != nil {
		return false, parseErr
	}
	return node.Name != nil, nil
}

func HasImportKeywords(file string) (bool, error) {
	fileSet := token.NewFileSet()
	node, parseErr := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
	if parseErr != nil {
		return false, parseErr
	}
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.IMPORT {
			return true, nil
		}
	}
	return false, nil
}

// TrimFileTitle 去除文件中的package和import声明，返回剩余数据
func TrimFileTitle(file string) (string, error) {
	// 解析文件以获取AST
	fileSet := token.NewFileSet()
	node, parseErr := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
	if parseErr != nil {
		return "", parseErr
	}

	// 创建一个新的AST打印配置
	var buf bytes.Buffer
	printConfig := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}

	// 遍历文件的顶层声明
	for i, decl := range node.Decls {
		// 跳过package和import声明
		genDecl, isGenDecl := decl.(*ast.GenDecl)
		if isGenDecl && (genDecl.Tok == token.IMPORT || genDecl.Tok == token.PACKAGE) {
			continue
		}
		// 打印其他节点
		err := printConfig.Fprint(&buf, fileSet, decl)
		if err != nil {
			return "", err
		}
		// 在声明之间添加换行符
		if i < len(node.Decls)-1 {
			buf.WriteString("\n\n") // 添加两个换行符以分隔顶层声明
		}
	}

	// 返回修改后的文件内容
	return buf.String(), nil
}

// FindFunctionInFile 在指定的文件中查找函数，注意，这里只查找非方法函数
func FindFunctionInFile(file string, functionName string) (*ast.FuncDecl, bool, error) {
	// 解析文件以获取AST
	fileSet := token.NewFileSet()
	node, parseErr := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
	if parseErr != nil {
		return nil, false, parseErr
	}
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Recv == nil {
				if funcDecl.Name.Name == functionName {
					return funcDecl, true, nil
				}
			}
		}
	}
	return nil, false, nil
}

// FindMethodInFile 在指定的文件中查找方法, 注意，这里只查找指定接收者类型的方法
func FindMethodInFile(file string, receiverTypeName, methodName string) (*ast.FuncDecl, bool, error) {
	// 解析文件以获取AST
	fileSet := token.NewFileSet()
	node, parseErr := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
	if parseErr != nil {
		return nil, false, parseErr
	}
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil {
				for _, recv := range funcDecl.Recv.List {
					if recvType, ok := recv.Type.(*ast.StarExpr); ok {
						if ident, ok := recvType.X.(*ast.Ident); ok && ident.Name == receiverTypeName {
							if funcDecl.Name.Name == methodName {
								return funcDecl, true, nil
							}
						}
					} else if ident, ok := recv.Type.(*ast.Ident); ok {
						if ident.Name == receiverTypeName {
							if funcDecl.Name.Name == methodName {
								return funcDecl, true, nil
							}
						}
					}
				}
			}
		}
	}
	return nil, false, nil
}
