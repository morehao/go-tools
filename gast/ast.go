package gast

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
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

// AddMethodToInterfaceInFile 在指定的文件中查找指定的接口，并在其中添加指定的方法
func AddMethodToInterfaceInFile(file, interfaceName, receiverTypeName, methodName string) error {

	// 查找接收者类型的方法声明
	methodDecl, found, findErr := FindMethodInFile(file, receiverTypeName, methodName)
	if findErr != nil {
		return findErr
	}
	if !found {
		return errors.New("method not found: " + methodName)
	}
	methodType := methodDecl.Type

	methodExists := false

	src, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	fileSet := token.NewFileSet()
	node, parseFileErr := parser.ParseFile(fileSet, file, src, parser.ParseComments)
	if parseFileErr != nil {
		return parseFileErr
	}

	// 遍历文件中的所有类型声明，查找指定的接口
	ast.Inspect(node, func(n ast.Node) bool {
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
		if err := format.Node(&buf, fileSet, node); err != nil {
			return err
		}

		// 重写文件
		if err := os.WriteFile(file, buf.Bytes(), 0644); err != nil {
			return err
		}
	}

	return nil
}

// AddContentToFunc 将指定的内容添加到指定文件中的指定函数的末尾
func AddContentToFunc(content, functionName, functionFilepath string) error {
	// 解析整个文件
	fileSet := token.NewFileSet()
	node, parseErr := parser.ParseFile(fileSet, functionFilepath, nil, parser.AllErrors)
	if parseErr != nil {
		return fmt.Errorf("failed to parse file: %w", parseErr)
	}

	// 查找目标函数
	var funcDecl *ast.FuncDecl
	ast.Inspect(node, func(n ast.Node) bool {
		if f, ok := n.(*ast.FuncDecl); ok && f.Name.Name == functionName {
			funcDecl = f
			return false
		}
		return true
	})
	if funcDecl == nil {
		return errors.New("function does not exist")
	}

	// 直接插入内容
	newStmt := &ast.ExprStmt{
		X: &ast.BasicLit{
			Kind:  token.STRING,
			Value: content,
		},
	}
	funcDecl.Body.List = append(funcDecl.Body.List, newStmt)

	// 使用 bytes.Buffer 处理修改后的内容
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fileSet, node); err != nil {
		return fmt.Errorf("failed to write updated content: %w", err)
	}

	// 打开已存在的目标文件
	file, openErr := os.OpenFile(functionFilepath, os.O_WRONLY|os.O_TRUNC, 0644)
	if openErr != nil {
		return fmt.Errorf("failed to open destination file: %v", openErr)
	}
	defer file.Close()

	// 将处理后的代码写回文件
	if _, writeErr := file.Write(buf.Bytes()); writeErr != nil {
		return fmt.Errorf("failed to write content: %w", writeErr)
	}
	return nil
}

// AddFunction 将指定的函数内容添加到指定文件中，如果文件不存在包声明则添加包声明
func AddFunction(content, functionFilepath, pkgName string) error {
	// 解析目标文件
	fileSet := token.NewFileSet()
	node, parseErr := parser.ParseFile(fileSet, functionFilepath, nil, parser.ParseComments)
	if parseErr != nil {
		// 如果文件解析失败，认为文件不存在或是空文件，创建新的文件节点
		node = &ast.File{
			Name: &ast.Ident{Name: pkgName},
		}
	}

	// 解析新的函数声明
	newFuncNode, parseFuncErr := parser.ParseFile(fileSet, "", "package "+pkgName+"\n"+content, parser.ParseComments)
	if parseFuncErr != nil {
		return fmt.Errorf("failed to parse new function: %w", parseFuncErr)
	}

	// 查找新的函数声明
	var newFuncDecl *ast.FuncDecl
	for _, decl := range newFuncNode.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			newFuncDecl = funcDecl
			break
		}
	}
	if newFuncDecl == nil {
		return fmt.Errorf("no function declaration found in content")
	}

	// 检查目标文件中是否有包声明
	hasPackage := node.Name != nil && node.Name.Name != ""

	// 如果目标文件没有包声明，设置包声明
	if !hasPackage {
		node.Name = &ast.Ident{Name: pkgName}
	}

	// 将新的函数声明添加到目标文件的声明列表中
	node.Decls = append(node.Decls, newFuncDecl)

	// 使用 bytes.Buffer 处理修改后的内容
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fileSet, node); err != nil {
		return fmt.Errorf("failed to write updated content: %w", err)
	}

	// 打开目标文件进行写入
	file, openErr := os.OpenFile(functionFilepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if openErr != nil {
		return fmt.Errorf("failed to open destination file: %v", openErr)
	}
	defer file.Close()

	// 将处理后的代码写回文件
	if _, writeErr := file.Write(buf.Bytes()); writeErr != nil {
		return fmt.Errorf("failed to write content: %w", writeErr)
	}

	return nil
}
