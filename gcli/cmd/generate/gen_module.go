package generate

import (
	"fmt"
	"github.com/morehao/go-tools/codeGen"
	"github.com/morehao/go-tools/gast"
	"github.com/morehao/go-tools/gutils"
	"path/filepath"
	"strings"
	"text/template"
)

func genModule(workDir string) error {
	cfg := Cfg.CodeGen.Module
	tplDir := filepath.Join(workDir, cfg.TplDir)
	rootDir := filepath.Join(workDir, cfg.InternalAppRootDir)
	layerDirMap := map[codeGen.LayerName]string{
		codeGen.LayerNameErrorCode: filepath.Join(filepath.Dir(rootDir), "/pkg"),
	}
	analysisCfg := &codeGen.ModuleCfg{
		CommonConfig: codeGen.CommonConfig{
			TplDir:      tplDir,
			PackageName: cfg.PackageName,
			RootDir:     rootDir,
			LayerDirMap: layerDirMap,
			TplFuncMap: template.FuncMap{
				TplFuncIsSysField: IsSysField,
			},
		},
		TableName: cfg.TableName,
	}
	gen := codeGen.NewGenerator()
	analysisRes, analysisErr := gen.AnalysisModuleTpl(MysqlClient, analysisCfg)
	if analysisErr != nil {
		return fmt.Errorf("analysis module tpl error: %v", analysisErr)
	}

	var genParamsList []codeGen.GenParamsItem
	for _, v := range analysisRes.TplAnalysisList {
		var modelFields []ModelField
		for _, field := range v.ModelFields {
			modelFields = append(modelFields, ModelField{
				FieldName:          gutils.ReplaceIdToID(field.FieldName),
				FieldLowerCaseName: gutils.SnakeToLowerCamel(field.FieldName),
				FieldType:          field.FieldType,
				ColumnName:         field.ColumnName,
				ColumnType:         field.ColumnType,
				Comment:            field.Comment,
				IsPrimaryKey:       field.ColumnKey == codeGen.ColumnKeyPRI,
			})
		}

		genParamsList = append(genParamsList, codeGen.GenParamsItem{
			TargetDir:      v.TargetDir,
			TargetFileName: v.TargetFilename,
			Template:       v.Template,
			ExtraParams: ModuleExtraParams{
				ServiceName:            Cfg.CodeGen.ServiceName,
				PackageName:            analysisRes.PackageName,
				PackagePascalName:      analysisRes.PackagePascalName,
				ProjectRootDir:         cfg.ProjectRootDir,
				TableName:              analysisRes.TableName,
				Description:            cfg.Description,
				StructName:             analysisRes.StructName,
				ReceiverTypeName:       gutils.FirstLetterToLower(analysisRes.StructName),
				ReceiverTypePascalName: analysisRes.StructName,
				ApiDocTag:              cfg.ApiDocTag,
				ApiGroup:               cfg.ApiGroup,
				ApiPrefix:              strings.TrimSuffix(cfg.ApiPrefix, "/"),
				Template:               v.Template,
				ModelFields:            modelFields,
			},
		})

	}
	genParams := &codeGen.GenParams{
		ParamsList: genParamsList,
	}
	if err := gen.Gen(genParams); err != nil {
		return err
	}

	// 注册路由
	routerCallContent := fmt.Sprintf("%sRouter(routerGroup)", gutils.FirstLetterToLower(analysisRes.StructName))
	routerEnterFilepath := filepath.Join(rootDir, "/router/enter.go")
	if err := gast.AddContentToFunc(routerEnterFilepath, "RegisterRouter", routerCallContent); err != nil {
		return fmt.Errorf("appendContentToFunc error: %v", err)
	}
	return nil
}

type ModuleExtraParams struct {
	ServiceName            string
	ProjectRootDir         string
	PackageName            string
	PackagePascalName      string
	TableName              string
	Description            string
	StructName             string
	ReceiverTypeName       string
	ReceiverTypePascalName string
	ApiGroup               string
	ApiPrefix              string
	ApiDocTag              string
	Template               *template.Template
	ModelFields            []ModelField
}
