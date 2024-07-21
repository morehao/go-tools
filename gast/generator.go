package gast

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

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

// AddMethodToInterfaceInFileV2 adds a method from a receiver type to an interface
func AddMethodToInterfaceInFileV2(filePath, receiverName, methodName, interfaceName string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	var methodDecl *ast.FuncDecl

	// Find the method declaration
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				if starExpr, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == receiverName && fn.Name.Name == methodName {
						methodDecl = fn
						return false
					}
				}
			}
		}
		return true
	})

	if methodDecl == nil {
		return fmt.Errorf("method %s not found in receiver %s", methodName, receiverName)
	}

	// Extract method signature without comments
	methodSig := removeCommentsFromFuncType(methodDecl.Type)

	// Find the interface and add method signature
	var interfaceType *ast.InterfaceType

	ast.Inspect(node, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if it, ok := ts.Type.(*ast.InterfaceType); ok && ts.Name.Name == interfaceName {
				interfaceType = it
				return false
			}
		}
		return true
	})

	if interfaceType == nil {
		return fmt.Errorf("interface %s not found", interfaceName)
	}

	// Add method signature to the interface
	methodField := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(methodName)},
		Type:  methodSig,
	}
	interfaceType.Methods.List = append(interfaceType.Methods.List, methodField)

	// Write modified AST back to file
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if err := printer.Fprint(f, fset, node); err != nil {
		return fmt.Errorf("failed to print AST: %w", err)
	}

	return nil
}

// removeCommentsFromFuncType removes comments from function type
func removeCommentsFromFuncType(ft *ast.FuncType) *ast.FuncType {
	return &ast.FuncType{
		Params:  removeCommentsFromFieldList(ft.Params),
		Results: removeCommentsFromFieldList(ft.Results),
	}
}

// removeCommentsFromFieldList removes comments from field list
func removeCommentsFromFieldList(fl *ast.FieldList) *ast.FieldList {
	if fl == nil {
		return nil
	}
	newList := make([]*ast.Field, len(fl.List))
	for i, f := range fl.List {
		newList[i] = &ast.Field{
			Names: f.Names,
			Type:  f.Type,
		}
	}
	return &ast.FieldList{List: newList}
}

// AddMethodToInterface 将指定接收者类型的方法添加到指定文件中的接口中。
func AddMethodToInterface(filePath, receiverType, methodName, interfaceName string) error {
	// 获取方法声明字符串。
	methodDecl, err := getMethodDeclaration(filePath, receiverType, methodName)
	if err != nil {
		return err
	}

	// 将方法声明添加到接口中。
	err = addContentToInterfaceInFile(filePath, methodDecl, interfaceName)
	if err != nil {
		return err
	}

	return nil
}

func getMethodDeclaration(filePath, receiverType, methodName string) (string, error) {
	// 解析文件以获取 AST。
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)
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

// addContentToInterfaceInFile 将给定内容添加到文件中指定的接口中。
func addContentToInterfaceInFile(filePath, content, interfaceName string) error {
	// 读取文件内容。
	lines, err := readLines(filePath)
	if err != nil {
		return err
	}

	// 在接口定义中插入内容。
	inserted, err := insertIntoInterface(lines, content, interfaceName)
	if err != nil {
		return err
	}
	if !inserted {
		return errors.New("接口未找到或内容已存在")
	}

	// 将修改后的内容写回文件。
	return writeLines(filePath, lines)
}

// readLines 读取文件的所有行。
func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// insertIntoInterface 在接口定义中插入内容。
func insertIntoInterface(lines []string, content, interfaceName string) (bool, error) {
	foundInterface := false
	inserted := false
	for i, line := range lines {
		if strings.Contains(line, fmt.Sprintf("type %s interface {", interfaceName)) {
			foundInterface = true
		} else if foundInterface && strings.Contains(line, "}") {
			lines[i] = "\t" + content + "\n" + line
			inserted = true
			break
		}
	}
	return inserted, nil
}

// writeLines 将所有行写回文件。
func writeLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

// exprToString 将 AST 表达式转换为字符串，只返回类型名称。
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident: // 标识符
		return t.Name
	case *ast.SelectorExpr: // 选择表达式，如包名.类型名
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr: // 指针类型
		return "*" + exprToString(t.X)
	case *ast.ArrayType: // 数组类型
		return "[]" + exprToString(t.Elt)
	// 这里可以添加更多的类型处理，如 Map, Chan, Func 等。
	default:
		return ""
	}
}
