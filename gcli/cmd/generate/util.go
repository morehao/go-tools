package generate

import (
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	TplFuncIsSysField = "isSysField"
)

func IsSysField(name string) bool {
	sysFieldMap := map[string]struct{}{
		"Id":        {},
		"CreatedAt": {},
		"CreatedBy": {},
		"UpdatedAt": {},
		"UpdatedBy": {},
		"DeletedAt": {},
		"DeletedBy": {},
	}
	_, ok := sysFieldMap[name]
	return ok
}

// CopyEmbeddedTemplatesToTempDir 将嵌入的模板文件复制到临时目录，并返回该目录的路径。
func CopyEmbeddedTemplatesToTempDir(embeddedFS embed.FS, root string) (string, error) {
	// 创建一个临时目录来存放模板文件
	tempDir, err := os.MkdirTemp("", "codegen_templates")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// 将嵌入的模板文件复制到临时目录
	err = fs.WalkDir(embeddedFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			data, readErr := embeddedFS.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			// 保持目录结构
			relPath, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			targetPath := filepath.Join(tempDir, relPath)
			if mkDirErr := os.MkdirAll(filepath.Dir(targetPath), 0755); mkDirErr != nil {
				return mkDirErr
			}
			if writeErr := ioutil.WriteFile(targetPath, data, 0644); writeErr != nil {
				return writeErr
			}
		}
		return nil
	})
	if err != nil {
		// 如果复制失败，清理临时目录
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to copy templates: %v", err)
	}

	return tempDir, nil
}
