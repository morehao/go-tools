package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/morehao/go-tools/codeGen"
	"github.com/morehao/go-tools/gast"
	"github.com/morehao/go-tools/gutils"
)

func genApi() error {
	apiGenCfg := cfg.CodeGen.Api

	// 使用工具函数复制嵌入的模板文件到临时目录
	tplDir, err := CopyEmbeddedTemplatesToTempDir(templatesFS, "template/api")
	if err != nil {
		return err
	}
	// 清理临时目录
	defer os.RemoveAll(tplDir)

	rootDir := filepath.Join(workDir, apiGenCfg.InternalAppRootDir)
	layerDirMap := map[codeGen.LayerName]string{
		codeGen.LayerNameErrorCode: filepath.Join(filepath.Dir(rootDir), "/pkg"),
	}
	analysisCfg := &codeGen.ApiCfg{
		CommonConfig: codeGen.CommonConfig{
			TplDir:      tplDir,
			PackageName: apiGenCfg.PackageName,
			RootDir:     rootDir,
			LayerDirMap: layerDirMap,
		},
		TargetFilename: apiGenCfg.TargetFilename,
	}
	gen := codeGen.NewGenerator()
	analysisRes, analysisErr := gen.AnalysisApiTpl(analysisCfg)
	if analysisErr != nil {
		return fmt.Errorf("analysis api tpl error: %v", analysisErr)
	}
	receiverTypePascalName := gutils.SnakeToPascal(apiGenCfg.SubModuleName)
	receiverTypeName := gutils.FirstLetterToLower(receiverTypePascalName)
	var genParamsList []codeGen.GenParamsItem
	var isNewRouter, isNewController bool
	var controllerFilepath, serviceFilepath string
	for _, v := range analysisRes.TplAnalysisList {
		switch v.LayerName {
		case codeGen.LayerNameRouter:
			if v.TargetFileExist {
				goFilepath := filepath.Join(v.TargetDir, v.TargetFilename)
				funcName := fmt.Sprintf("%sRouter", gutils.FirstLetterToLower(apiGenCfg.SubModuleName))
				_, hasFunc, findFuncErr := gast.FindFunction(goFilepath, funcName)
				if findFuncErr != nil {
					return fmt.Errorf("find function error: %v", findFuncErr)
				}
				isNewRouter = !hasFunc
			} else {
				isNewRouter = true
			}
		case codeGen.LayerNameController:
			controllerFilepath = filepath.Join(v.TargetDir, v.TargetFilename)
			isNewController = !v.TargetFileExist
		case codeGen.LayerNameService:
			serviceFilepath = filepath.Join(v.TargetDir, v.TargetFilename)
		}

		genParamsList = append(genParamsList, codeGen.GenParamsItem{
			TargetDir:      v.TargetDir,
			TargetFileName: v.TargetFilename,
			Template:       v.Template,
			ExtraParams: ApiExtraParams{
				PackageName:            analysisRes.PackageName,
				PackagePascalName:      analysisRes.PackagePascalName,
				ProjectRootDir:         apiGenCfg.ProjectRootDir,
				TargetFileExist:        v.TargetFileExist,
				IsNewRouter:            isNewRouter,
				Description:            apiGenCfg.Description,
				ReceiverTypeName:       receiverTypeName,
				ReceiverTypePascalName: receiverTypePascalName,
				HttpMethod:             apiGenCfg.HttpMethod,
				FunctionName:           gutils.FirstLetterToUpper(apiGenCfg.FunctionName),
				ApiDocTag:              apiGenCfg.ApiDocTag,
				ApiPrefix:              strings.TrimSuffix(apiGenCfg.ApiPrefix, "/"),
				ApiSuffix:              strings.TrimLeft(apiGenCfg.ApiSuffix, "/"),
				ApiGroup:               apiGenCfg.ApiGroup,
				Template:               v.Template,
			},
		})

	}
	genParams := &codeGen.GenParams{
		ParamsList: genParamsList,
	}
	if err := gen.Gen(genParams); err != nil {
		return err
	}

	if !isNewController {
		// 将方法添加到interface接口中
		controllerInterfaceName := fmt.Sprintf("%sCtr", receiverTypePascalName)
		if err := gast.AddMethodToInterface(controllerFilepath, receiverTypeName+"Ctr", apiGenCfg.FunctionName, controllerInterfaceName); err != nil {
			return fmt.Errorf("add controller method to interface error: %w", err)
		}
		serviceInterfaceName := fmt.Sprintf("%sSvc", receiverTypePascalName)
		if err := gast.AddMethodToInterface(serviceFilepath, receiverTypeName+"Svc", apiGenCfg.FunctionName, serviceInterfaceName); err != nil {
			return fmt.Errorf("add service method to interface error: %w", err)
		}
	}

	// 	注册路由
	if isNewRouter {
		routerCallContent := fmt.Sprintf("%sRouter(routerGroup)", receiverTypeName)
		routerEnterFilepath := filepath.Join(rootDir, "/router/enter.go")
		if err := gast.AddContentToFunc(routerEnterFilepath, "RegisterRouter", routerCallContent); err != nil {
			return fmt.Errorf("new router appendContentToFunc error: %v", err)
		}
	} else {
		routerCallContent := fmt.Sprintf(`routerGroup.%s("/%s", %sCtr.%s) // %s`, apiGenCfg.HttpMethod, apiGenCfg.ApiSuffix, receiverTypeName, apiGenCfg.FunctionName, apiGenCfg.Description)
		routerEnterFilepath := filepath.Join(rootDir, fmt.Sprintf("/router/%s.go", gutils.TrimFileExtension(apiGenCfg.TargetFilename)))
		if err := gast.AddContentToFuncWithLineNumber(routerEnterFilepath, fmt.Sprintf("%sRouter", receiverTypeName), routerCallContent, -2); err != nil {
			return fmt.Errorf("appendContentToFunc error: %v", err)
		}
	}
	return nil
}

type ApiExtraParams struct {
	ServiceName            string
	ProjectRootDir         string
	PackageName            string
	PackagePascalName      string
	Description            string
	TargetFileExist        bool
	IsNewRouter            bool
	HttpMethod             string
	FunctionName           string
	ReceiverTypeName       string
	ReceiverTypePascalName string
	ApiGroup               string
	ApiPrefix              string
	ApiSuffix              string
	ApiDocTag              string
	Template               *template.Template
}
