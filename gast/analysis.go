package gast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
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

// FindFunction 在指定的文件中查找函数，注意，这里只查找非方法函数
func FindFunction(file string, functionName string) (*ast.FuncDecl, bool, error) {
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

// FindMethod 在指定的文件中查找方法, 注意，这里只查找指定接收者类型的方法
func FindMethod(file string, receiverTypeName, methodName string) (*ast.FuncDecl, bool, error) {
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

// interfaceContainsMethod 检查指定接口是否已经包含了指定的方法。
func interfaceContainsMethod(filePath, interfaceName, methodName string) (bool, error) {
	// 解析文件以获取 AST。
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, filePath, nil, 0)
	if err != nil {
		return false, err
	}

	// 遍历 AST 查找指定的接口。
	var contains bool
	ast.Inspect(node, func(n ast.Node) bool {
		// 查找接口类型声明。
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true // 继续遍历 AST。
		}
		if typeSpec.Name.Name == interfaceName {
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return true
			}
			// 遍历接口中的方法，检查是否包含指定的方法。
			for _, method := range interfaceType.Methods.List {
				if len(method.Names) > 0 && method.Names[0].Name == methodName {
					contains = true
					return false // 停止遍历 AST。
				}
			}
		}
		return true
	})

	return contains, nil
}

// getMethodDeclaration 获取指定方法声明
func getMethodDeclaration(filePath, receiverType, methodName string) (string, error) {
	// 解析文件以获取 AST。
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, filePath, nil, 0)
	if err != nil {
		return "", err
	}

	// 遍历 AST 查找指定的方法。
	var methodDecl string
	ast.Inspect(node, func(n ast.Node) bool {
		// 查找函数声明。
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true // 继续遍历 AST。
		}

		// 检查是否为指定的接收者和方法。
		if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
			recvType, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
			if !ok {
				return true
			}
			ident, ok := recvType.X.(*ast.Ident)
			if !ok || ident.Name != receiverType {
				return true
			}
			if funcDecl.Name.Name == methodName {
				// 找到方法，构建声明字符串。
				methodDecl = methodName + fieldListToString(funcDecl.Type.Params, false) + " " + fieldListToString(funcDecl.Type.Results, true)
				return false // 停止遍历 AST。
			}
		}
		return true
	})

	if methodDecl == "" {
		return "", fmt.Errorf("未找到接收者类型 '%s' 和方法名 '%s' 的方法声明", receiverType, methodName)
	}

	return methodDecl, nil
}

// fieldListToString 将 *ast.FieldList 转换为字符串表示，用于参数和返回值。
// isResults 参数指示这个字段列表是否是函数的返回值列表。
func fieldListToString(fl *ast.FieldList, isResults bool) string {
	if fl == nil || len(fl.List) == 0 {
		if isResults {
			return ""
		}
		return "()"
	}

	var fields []string
	for _, field := range fl.List {
		typeStr := exprToString(field.Type)
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				fields = append(fields, name.Name+" "+typeStr)
			}
		} else {
			fields = append(fields, typeStr)
		}
	}

	// 如果是返回值列表且有多个字段，或者是参数列表，用括号括起来。
	if isResults && len(fields) > 1 {
		return "(" + strings.Join(fields, ", ") + ")"
	} else if isResults && len(fields) == 1 {
		// 单个返回值不需要括号
		return strings.Join(fields, ", ")
	}
	// 参数列表始终需要括号
	return "(" + strings.Join(fields, ", ") + ")"
}
