package codeGen

import (
	"fmt"
	"github.com/morehao/go-tools/gutils"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const (
	tplFileSuffix = ".tpl"
	goFileSuffix  = ".go"

	tplLayerNameRouter     = "router"
	tplLayerNameController = "controller"
	tplLayerNameService    = "service"
	tplLayerNameDto        = "dto"
	tplLayerNameRequest    = "request"
	tplLayerNameResponse   = "response"
	tplLayerNameModel      = "model"

	tplLayerPrefixController = "ctr"
	tplLayerPrefixService    = "svc"
	tplLayerPrefixDto        = "dto"
	tplLayerPrefixModel      = "dao"
)

var layerPrefixMap = map[string]string{
	tplLayerNameController: tplLayerPrefixController,
	tplLayerNameService:    tplLayerPrefixService,
	tplLayerNameModel:      tplLayerPrefixModel,
	tplLayerNameDto:        tplLayerPrefixDto,
}

var layerSpecialNameMap = map[string]string{
	tplLayerNameRequest:  tplLayerPrefixDto,
	tplLayerNameResponse: tplLayerPrefixDto,
}

var layerAppendDirMap = map[string]string{
	tplLayerNameRouter: "routerHttp",
}

type tplFile struct {
	filepath       string
	filename       string
	originFilename string
	layerName      string
	layerPrefix    string
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
	Template       *template.Template
	Filepath       string
	Filename       string
	OriginFilename string
	TargetFileName string
	TargetDir      string
	LayerName      string
	LayerPrefix    string
}

type ModelTemplateParamsItem struct {
	TemplateParamsItemBase
	ModelFields []ModelField
}

type ControllerTemplateParams struct {
	PackageName       string
	PackagePascalName string
	TemplateList      []TemplateParamsItemBase
}

type GenParam struct {
	Params []GenParamsItem
}

type GenParamsItem struct {
	Template       *template.Template
	TargetDir      string
	TargetFileName string
	Param          interface{}
}

// 获取指定目录下所有的模板文件
func getTplFiles(path string) ([]tplFile, error) {
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
		if gutils.GetFileSuffix(name) == tplFileSuffix {
			layerName := strings.TrimSuffix(name, fmt.Sprintf("%s%s", goFileSuffix, tplFileSuffix))
			if specialName, ok := layerSpecialNameMap[layerName]; ok {
				layerName = specialName
			}
			files = append(files, tplFile{
				filepath:       fmt.Sprintf("%s/%s", path, name),
				filename:       name,
				originFilename: name[:len(name)-len(tplFileSuffix)],
				layerName:      layerName,
				layerPrefix:    layerPrefixMap[layerName],
			})
		}
	}
	return files, nil
}

func buildTplCfg(tplFiles []tplFile, defaultFilename string) ([]tplCfg, error) {
	var tplList []tplCfg
	for _, file := range tplFiles {
		targetFileName := file.originFilename
		if file.layerName != tplLayerNameDto {
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
		otherContent, trimErr := trimFileTitle(tmpFilepath)
		if trimErr != nil {
			return trimErr
		}
		f, openErr := os.OpenFile(codeFilepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if openErr != nil {
			return openErr
		}
		_, writeErr := f.WriteString(otherContent)
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

func trimFileTitle(file string) (string, error) {
	content, readErr := os.ReadFile(file)
	if readErr != nil {
		return "", readErr
	}
	fileContent := string(content)

	// 正则表达式匹配 package 语句
	packagePattern := regexp.MustCompile(`package\s+\w+\s*\n`)
	// 查找 package 语句的位置
	// 查找 package 语句的位置
	packageMatch := packagePattern.FindStringIndex(fileContent)
	var importStartIndex int
	if packageMatch != nil {
		importStartIndex = packageMatch[1]
	} else {
		importStartIndex = 0
	}

	// 正则表达式匹配 import 块和单个 import 语句
	importBlockPattern := regexp.MustCompile(`(?s)import \((.|\n)*?\)\n`)
	singleImportPattern := regexp.MustCompile(`import ".*?"\n`)
	// 查找 import 块和单个 import 语句的位置
	importBlockMatch := importBlockPattern.FindStringIndex(fileContent[importStartIndex:])
	singleImportMatch := singleImportPattern.FindStringIndex(fileContent[importStartIndex:])

	// 确定 import 语句及其块的结束位置
	var importEndIndex int
	if importBlockMatch != nil {
		importEndIndex = importStartIndex + importBlockMatch[1]
	} else if singleImportMatch != nil {
		importEndIndex = importStartIndex + singleImportMatch[1]
	} else {
		importEndIndex = importStartIndex
	}

	// 分割文件内容
	otherContent := fileContent[importEndIndex:]
	return otherContent, nil
}
