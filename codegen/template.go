package codegen

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/morehao/golib/gast"
	"github.com/morehao/golib/gutils"
)

const (
	tplFileExtension = ".tpl"
	goFileExtension  = ".go"
)

type ModuleTplAnalysisRes struct {
	PackageName     string
	TableName       string
	StructName      string
	TplAnalysisList []ModuleTplAnalysisItem
}

type TplAnalysisItem struct {
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

type ModuleTplAnalysisItem struct {
	TplAnalysisItem
	ModelFields []ModelField
}

type ApiTplAnalysisRes struct {
	PackageName     string
	TplAnalysisList []TplAnalysisItem
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
func analysisTplFiles(cfg CommonConfig, defaultTargetFilename string) ([]TplAnalysisItem, error) {
	// 打开指定目录
	file, openErr := os.Open(cfg.TplDir)
	if openErr != nil {
		return nil, openErr
	}
	// 读取目录下所有文件
	tplFilenameList, readErr := file.Readdirnames(-1)
	if readErr != nil {
		return nil, readErr
	}
	var analysisList []TplAnalysisItem
	rootDir := cfg.RootDir
	for _, tplFilename := range tplFilenameList {
		// 判断是否是模板文件
		if gutils.GetFileExtension(tplFilename) != tplFileExtension {
			continue
		}

		// 构造文件层级名称，如controller
		defaultLayerName := LayerName(strings.TrimSuffix(tplFilename, fmt.Sprintf("%s%s", goFileExtension, tplFileExtension)))
		layerName := defaultLayerName
		if specialName, ok := defaultLayerSpecialNameMap[layerName]; ok {
			layerName = specialName
		}
		if customLayerName, ok := cfg.LayerNameMap[defaultLayerName]; ok {
			layerName = customLayerName
		}

		// 构造生成文件所在模块的模块前缀
		defaultLayerPrefix := defaultLayerPrefixMap[defaultLayerName]
		layerPrefix := defaultLayerPrefix
		if prefix, ok := cfg.LayerPrefixMap[defaultLayerName]; ok {
			layerPrefix = prefix
		}

		layerParentDir := rootDir
		// 构造生成文件所在目录的名称
		if customLayerParentDir, ok := cfg.LayerParentDirMap[defaultLayerName]; ok {
			layerParentDir = path.Join(layerParentDir, customLayerParentDir)
		}
		var targetDir string
		if defaultLayerPrefix.String() == "" {
			targetDir = path.Join(layerParentDir, string(layerName))
		} else {
			targetFileParentDir := fmt.Sprintf("%s%s", layerPrefix, strings.ToLower(gutils.SnakeToPascal(cfg.PackageName)))
			targetDir = path.Join(layerParentDir, string(layerName), targetFileParentDir)
		}

		// 构造生成文件的文件名称
		originFilename := gutils.TrimFileExtension(gutils.TrimFileExtension(tplFilename))
		var targetFilename string
		switch defaultLayerName {
		case LayerNameRequest, LayerNameResponse:
			targetFilename = fmt.Sprintf("%s%s", originFilename, goFileExtension)
		case LayerNameRouter, LayerNameCode:
			targetFilename = fmt.Sprintf("%s%s", gutils.CamelToSnakeCase(cfg.PackageName), goFileExtension)
		default:
			targetFilename = fmt.Sprintf("%s%s", gutils.TrimFileExtension(defaultTargetFilename), goFileExtension)
		}

		var targetFileExist bool
		if gutils.FileExists(path.Join(targetDir, targetFilename)) {
			targetFileExist = true
		}
		tplFilepath := path.Join(cfg.TplDir, tplFilename)
		tplInst := template.New(tplFilename).Funcs(cfg.TplFuncMap)
		fileTemplate, parseErr := tplInst.ParseFiles(tplFilepath)
		if parseErr != nil {
			return nil, parseErr
		}

		analysisList = append(analysisList, TplAnalysisItem{
			Template:        fileTemplate,
			TplFilepath:     tplFilepath,
			TplFilename:     tplFilename,
			LayerName:       layerName,
			LayerPrefix:     layerPrefix,
			OriginFilename:  originFilename,
			TargetDir:       targetDir,
			TargetFilename:  targetFilename,
			TargetFileExist: targetFileExist,
		})

	}
	return analysisList, nil
}

func createFile(targetDir, targetFilename string, tpl *template.Template, tplParam interface{}) error {
	if err := gutils.CreateDir(targetDir); err != nil {
		return err
	}
	codeFilepath := path.Join(targetDir, targetFilename)
	// 判断文件是否存在
	if exist := gutils.FileExists(codeFilepath); exist {
		// 如果存在，先写入一个临时文件，再对既有文件进行追加
		tempDir := path.Join(targetDir, "tmp")
		tmpFilepath := path.Join(tempDir, targetFilename)
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
			appendContent = "\n" + trimTitleContent
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
