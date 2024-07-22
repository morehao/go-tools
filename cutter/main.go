package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 检查参数
	if len(os.Args) < 3 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("介绍: 初始化项目")
		fmt.Println("用法: cutter target_path project_name")
		fmt.Println("说明: target_path: 目标项目路径, project_name: 项目名称")
		fmt.Println("示例: cutter /root/new-projects my_app")
		return
	}

	targetPath := os.Args[1]
	projectName := os.Args[2]

	// 获取当前执行目录
	currentDir := getCurrentDirectory()
	// 获取原始项目名称
	originalProjectName := filepath.Base(currentDir)

	// 检查原始项目是否为Go项目
	if !isGoProject(currentDir) {
		fmt.Println("原始项目不是Go项目，缺少go.mod文件")
		return
	}

	// 检查目标路径是否有效
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		fmt.Println("目标路径不存在:", targetPath)
		return
	}

	// 读取.gitignore文件并创建排除列表
	excludes, err := readGitignore(filepath.Join(currentDir, ".gitignore"))
	if err != nil {
		fmt.Println("读取.gitignore文件失败:", err)
		return
	}

	// 创建新项目目录
	newProjectPath := filepath.Join(targetPath, projectName)
	err = os.MkdirAll(newProjectPath, os.ModePerm)
	if err != nil {
		fmt.Println("创建新项目目录失败:", err)
		return
	}
	fmt.Println("新项目目录:", newProjectPath)

	// 复制文件
	err = copyDir(currentDir, newProjectPath, excludes)
	if err != nil {
		fmt.Println("复制文件失败:", err)
		return
	}

	// 替换文本
	err = filepath.Walk(newProjectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			read, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			newContents := strings.Replace(string(read), originalProjectName, projectName, -1)

			err = os.WriteFile(path, []byte(newContents), info.Mode())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("替换项目名称失败:", err)
		return
	}

	fmt.Println("成功生成项目，请移步", newProjectPath, "查看 😊")
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return dir
}

func isGoProject(path string) bool {
	_, err := os.Stat(filepath.Join(path, "go.mod"))
	return !os.IsNotExist(err)
}

func readGitignore(filename string) ([]string, error) {
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

func copyDir(src string, dst string, excludes []string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
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
		targetPath := strings.Replace(path, src, dst, 1)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		// 复制文件
		if !info.IsDir() {
			return copyFile(path, targetPath)
		}
		return nil
	})
}

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
