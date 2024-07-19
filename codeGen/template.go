package codeGen

import (
	"fmt"
	"github.com/morehao/go-tools/gast"
	"github.com/morehao/go-tools/gutils"
	"os"
	"strings"
	"text/template"
)

const (
	tplFileSuffix = ".tpl"
	goFileSuffix  = ".go"
)

type tplFile struct {
	filepath       string
	filename       string
	originFilename string
	layerName      LayerName
	layerPrefix    LayerPrefix
}

type tplCfg struct {
	template *template.Template
	tplFile
	targetFileName string
}

func (t *tplCfg) BuildTargetDir(rootDir, packageName string) string {
	if t.layerPrefix == "" {
		return fmt.Sprintf("%s/%s", rootDir, t.layerName)
	}
	layerDirName := fmt.Sprintf("%s%s", t.layerPrefix, gutils.SnakeToPascal(packageName))
	return fmt.Sprintf("%s/%s/%s", rootDir, t.layerName, layerDirName)
}

type ModelTemplateParamsRes struct {
	PackageName       string
	TableName         string
	PackagePascalName string
	StructName        string
	TemplateList      []ModelTemplateParamsItem
}

type TemplateParamsItemBase struct {
	Template        *template.Template
	TplFilepath     string
	TplFilename     string
	OriginFilename  string
	TargetDir       string
	TargetFilename  string
	TargetFileExist bool
	LayerName       LayerName
	LayerPrefix     LayerPrefix
}

type ModelTemplateParamsItem struct {
	TemplateParamsItemBase
	ModelFields []ModelField
}

type ControllerTemplateParamsRes struct {
	PackageName       string
	PackagePascalName string
	TemplateList      []TemplateParamsItemBase
}

type GenParams struct {
	ParamsList []GenParamsItem
}

type GenParamsItem struct {
	Template       *template.Template
	TargetDir      string
	TargetFileName string
	ExtraParams    interface{}
}

// 获取指定目录下所有的模板文件
func getTplFiles(tplDir string, layerNameMap map[LayerName]LayerName, layerPrefixMap map[LayerName]LayerPrefix) ([]tplFile, error) {
	// 打开指定目录
	file, openErr := os.Open(tplDir)
	if openErr != nil {
		return nil, openErr
	}
	// 读取目录下所有文件
	names, readErr := file.Readdirnames(-1)
	if readErr != nil {
		return nil, readErr
	}
	var files []tplFile
	for _, name := range names {

		// 判断是否是模板文件
		if gutils.GetFileSuffix(name) == tplFileSuffix {
			layerName := LayerName(strings.TrimSuffix(name, fmt.Sprintf("%s%s", goFileSuffix, tplFileSuffix)))
			if specialName, ok := defaultLayerSpecialNameMap[layerName]; ok {
				layerName = specialName
			}
			layerPrefix := defaultLayerPrefixMap[layerName]
			if prefix, ok := layerPrefixMap[layerName]; ok {
				layerPrefix = prefix
			}

			if specialName, ok := layerNameMap[layerName]; ok {
				layerName = specialName
			}
			files = append(files, tplFile{
				filepath:       fmt.Sprintf("%s/%s", tplDir, name),
				filename:       name,
				originFilename: name[:len(name)-len(tplFileSuffix)],
				layerName:      layerName,
				layerPrefix:    layerPrefix,
			})
		}
	}
	return files, nil
}

func buildTplCfg(tplFiles []tplFile, defaultFilename string) ([]tplCfg, error) {
	var tplList []tplCfg
	for _, file := range tplFiles {
		targetFileName := file.originFilename
		if file.layerName != LayerNameDto {
			targetFileName = defaultFilename
		}
		tplItem := tplCfg{
			tplFile:        file,
			targetFileName: targetFileName,
		}
		tplList = append(tplList, tplItem)
	}
	for i, tplItem := range tplList {
		tpl, parseErr := template.ParseFiles(tplItem.filepath)
		if parseErr != nil {
			return nil, parseErr
		}
		tplList[i].template = tpl
	}
	return tplList, nil
}

func createFile(targetDir, targetFileName string, tpl *template.Template, tplParam interface{}) error {
	if err := gutils.CreateDir(targetDir); err != nil {
		return err
	}
	codeFilepath := fmt.Sprintf("%s/%s", targetDir, targetFileName)
	// 判断文件是否存在
	if exist := gutils.FileExists(codeFilepath); exist {
		// 如果存在，先写入一个临时文件，再对既有文件进行追加
		tempDir := fmt.Sprintf("%s/tmp", targetDir)
		tmpFilepath := fmt.Sprintf("%s/%s", tempDir, targetFileName)
		if err := gutils.CreateDir(tempDir); err != nil {
			return err
		}
		tempF, openTempErr := os.OpenFile(tmpFilepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if openTempErr != nil {
			return openTempErr
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				panic(err)
			}
		}()
		if err := tpl.Execute(tempF, tplParam); err != nil {
			return err
		}

		// 判断文件是否包含package和import关键字
		hasPackage, checkPackageErr := gast.HasPackageKeywords(tmpFilepath)
		if checkPackageErr != nil {
			return checkPackageErr
		}
		hasImport, checkImportErr := gast.HasImportKeywords(tmpFilepath)
		if checkImportErr != nil {
			return checkImportErr
		}

		// 获取追加的文件内容
		var appendContent string
		if hasPackage || hasImport {
			trimTitleContent, trimErr := gast.TrimFileTitle(tmpFilepath)
			if trimErr != nil {
				return trimErr
			}
			appendContent = trimTitleContent
		} else {
			content, readErr := os.ReadFile(tmpFilepath)
			if readErr != nil {
				return readErr
			}
			appendContent = string(content)
		}

		// 追加到原文件
		f, openErr := os.OpenFile(codeFilepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if openErr != nil {
			return openErr
		}
		_, writeErr := f.WriteString(appendContent)
		if writeErr != nil {
			return writeErr
		}
	} else {
		f, openErr := os.OpenFile(codeFilepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if openErr != nil {
			return openErr
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()
		if err := tpl.Execute(f, &tplParam); err != nil {
			return err
		}
	}
	return nil
}
