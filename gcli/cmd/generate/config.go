package generate

import (
	"github.com/morehao/go-tools/dbClient"
	"github.com/morehao/go-tools/glog"
)

type Config struct {
	Mysql   dbClient.MysqlConfig `yaml:"mysql"`
	Log     glog.LoggerConfig    `yaml:"log"`
	CodeGen CodeGen              `yaml:"code_gen"`
}

type CodeGen struct {
	Mode   string       `yaml:"mode"`   // 生成模式，支持：module、model、api
	Module ModuleConfig `yaml:"module"` // 模块生成配置
	Model  ModelConfig  `yaml:"model"`  // 模型生成配置
	Api    ApiConfig    `yaml:"api"`    // 控制器生成配置
}

type ModuleConfig struct {
	InternalAppRootDir string `yaml:"internal_app_root_dir"` // 项目内当前项目的根目录，如internal/genCode
	ProjectRootDir     string `yaml:"project_root_dir"`      // 项目根目录，如go-gin-web
	Description        string `yaml:"description"`           // 描述
	ApiDocTag          string `yaml:"api_doc_tag"`           // api文档tag
	ApiGroup           string `yaml:"api_group"`             // api分组
	ApiPrefix          string `yaml:"api_prefix"`            // api前缀
	PackageName        string `yaml:"package_name"`          // 包名
	TableName          string `yaml:"table_name"`            // 表名
}

type ModelConfig struct {
	InternalAppRootDir string `yaml:"internal_app_root_dir"` // 项目内当前项目的根目录，如internal/genCode
	ProjectRootDir     string `yaml:"project_root_dir"`      // import目录前缀
	Description        string `yaml:"description"`           // 描述
	PackageName        string `yaml:"package_name"`          // 包名
	TableName          string `yaml:"table_name"`            // 表名
}

type ApiConfig struct {
	InternalAppRootDir string `yaml:"internal_app_root_dir"` // 项目内当前项目的根目录，如internal/genCode
	ProjectRootDir     string `yaml:"project_root_dir"`      // import目录前缀
	Description        string `yaml:"description"`           // 描述
	PackageName        string `yaml:"package_name"`          // 包名，如user
	SubModuleName      string `yaml:"sub_module_name"`       // 子模块名，如userAccount
	TargetFilename     string `yaml:"target_filename"`       // 目标文件名，生成的代码写入的目标文件名
	HttpMethod         string `yaml:"http_method"`           // http方法
	FunctionName       string `yaml:"function_name"`         // 函数名
	ApiDocTag          string `yaml:"api_doc_tag"`           // api文档tag
	ApiGroup           string `yaml:"api_group"`             // api组，RouterGroup对应的接口路径
	ApiPrefix          string `yaml:"api_prefix"`            // api前缀
	ApiSuffix          string `yaml:"api_suffix"`            // api后缀
}
