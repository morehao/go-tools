package autoCode

import (
	"fmt"
	"github.com/morehao/go-tools/utils"
	"os"
	"strings"
	"text/template"
)

const (
	tplFileSuffix = ".tpl"
)

type tplFile struct {
	filepath       string
	filename       string
	originFilename string
	layerName      string
}

type tplCfg struct {
	template template.Template
	tplFile
	targetFilepath string
}

type tplParam struct {
	PackageName   string
	TableName     string
	PackagePascal string
	StructName    string
}

func Tmpl() {
	// 模板定义
	tepl := "My name is {{.}}"
	// 解析模板
	tmpl, err := template.New("test").Parse(tepl)
	if err != nil {
		panic(err)
	}
	// 渲染模板
	if err := tmpl.Execute(os.Stdout, "Jack"); err != nil {
		panic(err)
	}

}

// 获取指定目录下所有的模板文件
func getTmplFiles(path string) ([]tplFile, error) {
	// 打开指定目录
	dir, openErr := os.Open(path)
	if openErr != nil {
		return nil, openErr
	}
	// 读取目录下所有文件
	names, readErr := dir.Readdirnames(-1)
	if readErr != nil {
		return nil, readErr
	}
	var files []tplFile
	for _, name := range names {

		// 判断是否是模板文件
		if utils.GetFileSuffix(name) == tplFileSuffix {
			layerName := strings.TrimSuffix(name, fmt.Sprintf(".go%s", tplFileSuffix))
			files = append(files, tplFile{
				filepath:       fmt.Sprintf("%s/%s", path, name),
				filename:       name,
				originFilename: name[:len(name)-len(tplFileSuffix)],
				layerName:      layerName,
			})
		}
	}
	return files, nil
}

func createFile(packageName, tableName string, tplList []tplCfg) error {
	packagePascal := utils.SnakeToPascal(packageName)
	structName := utils.SnakeToPascal(tableName)
	tmplParam := tplParam{
		PackageName:   packageName,
		TableName:     tableName,
		PackagePascal: packagePascal,
		StructName:    structName,
	}
	tempPath := "./tempAutoCode"
	for _, tplItem := range tplList {
		f, err := os.Open(fmt.Sprintf("%s/%s", tempPath, tplItem.originFilename))
		if err != nil {
			return err
		}
		if err = tplItem.template.Execute(f, tmplParam); err != nil {
			return err
		}
		_ = f.Close()
	}
	return nil
}
