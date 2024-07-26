package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 解析命令行参数
	destination := flag.String("d", "", "指定新项目的目录")
	flag.Usage = func() {
		fmt.Println("Usage: gcutter -d <new_project_directory>")
	}

	flag.Parse()

	// 检查是否提供了目标目录
	if *destination == "" {
		fmt.Println("必须提供新项目的目录")
		flag.Usage()
		return
	}
	newProjectPath := *destination
	cutter(newProjectPath)
}

func cutter(newProjectPath string) {
	// 获取当前执行目录，确认它是Go项目
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前目录失败:", err)
		return
	}
	if !isGoProject(currentDir) {
		fmt.Println("当前目录不是Go项目")
		return
	}

	// 获取模板项目名称
	templateName := filepath.Base(currentDir)

	newProjectName := filepath.Base(newProjectPath)

	// 确认新项目目录不存在或为空
	if _, err := os.Stat(newProjectPath); !os.IsNotExist(err) {
		fmt.Println("新项目目录已存在:", newProjectPath)
		return
	}

	// 创建新项目目录
	if err := os.MkdirAll(newProjectPath, os.ModePerm); err != nil {
		fmt.Println("创建新项目目录失败:", err)
		return
	}

	// 读取.gitignore文件并创建排除列表
	excludes, err := readGitignore(filepath.Join(currentDir, ".gitignore"))
	if err != nil {
		fmt.Println("读取.gitignore文件失败:", err)
		return
	}

	// 复制模板项目到新项目目录，并替换import路径
	if err := copyAndReplace(currentDir, newProjectPath, templateName, newProjectName, excludes); err != nil {
		fmt.Println("复制模板项目时出错:", err)
		return
	}
	if err := removeGitDir(newProjectPath); err != nil {
		fmt.Println("删除.git目录失败:", err)
		return
	}
	fmt.Println("新项目创建成功:", newProjectPath)
}

// isGoProject 检查指定路径是否为Go项目（是否包含go.mod文件）
func isGoProject(path string) bool {
	_, err := os.Stat(filepath.Join(path, "go.mod"))
	return !os.IsNotExist(err)
}

// readGitignore 读取.gitignore文件并返回排除列表
func readGitignore(filename string) ([]string, error) {
	// 如果.gitignore文件不存在，则返回一个空列表
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []string{}, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

// copyAndReplace 复制模板项目到新项目目录，并替换import路径
func copyAndReplace(srcDir, dstDir, oldName, newName string, excludes []string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否排除
		for _, exclude := range excludes {
			if strings.Contains(path, exclude) {
				return nil
			}
		}

		// 创建目标目录
		targetPath := strings.Replace(path, srcDir, dstDir, 1)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		// 复制文件并替换 import 路径
		if strings.HasSuffix(info.Name(), ".go") {
			return copyAndReplaceGoFile(path, targetPath, oldName, newName)
		}

		// 复制其他文件
		return copyFile(path, targetPath)
	})
	if err != nil {
		return err
	}
	if err := modifyGoMod(dstDir, newName); err != nil {
		return err
	}
	return err
}

// copyAndReplaceGoFile 复制并替换 Go 文件中的 import 路径
func copyAndReplaceGoFile(srcFile, dstFile, oldName, newName string) error {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, srcFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// 遍历文件中的所有 import 语句，替换路径中的 oldName 为 newName
	ast.Inspect(node, func(n ast.Node) bool {
		importSpec, ok := n.(*ast.ImportSpec)
		if ok {
			importPath := strings.Trim(importSpec.Path.Value, `"`)
			if strings.Contains(importPath, oldName) {
				updatedImportPath := strings.Replace(importPath, oldName, newName, -1)
				importSpec.Path.Value = fmt.Sprintf(`"%s"`, updatedImportPath)
			}
		}
		return true
	})

	// 将更新后的代码写入目标文件
	file, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := format.Node(file, fs, node); err != nil {
		return err
	}
	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// 修改go.mod中的包名
func modifyGoMod(dstDir, moduleName string) error {
	// 读取go.mod文件
	modFilepath := filepath.Join(dstDir, "go.mod")
	content, err := os.ReadFile(modFilepath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	// 解析go.mod文件
	modFile, err := modfile.Parse(modFilepath, content, nil)
	if err != nil {
		return fmt.Errorf("failed to parse go.mod: %w", err)
	}

	// 修改模块名称
	if err := modFile.AddModuleStmt(moduleName); err != nil {
		return fmt.Errorf("failed to add module statement: %w", err)
	}

	// 将修改后的内容格式化回字节切片
	newContent, err := modFile.Format()
	if err != nil {
		return fmt.Errorf("failed to format new go.mod content: %w", err)
	}

	// 写入新的go.mod文件
	err = os.WriteFile(modFilepath, newContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write new go.mod: %w", err)
	}
	fmt.Println(fmt.Sprintf("Successfully modified go.mod to %s", moduleName))
	return nil
}

// 删除.git目录
// removeGitDir 删除指定目录下的.git文件夹
func removeGitDir(dstDir string) error {
	gitDir := filepath.Join(dstDir, ".git")
	err := os.RemoveAll(gitDir)
	if err != nil {
		return fmt.Errorf("failed to remove .git directory: %w", err)
	}
	fmt.Println("Successfully removed .git directory")
	return nil
}
