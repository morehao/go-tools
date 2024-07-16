package codeGen

type ModuleCfg struct {
	PackageName    string                    `validate:"required"` // 包名
	TableName      string                    `validate:"required"` // 表名
	ColumnTypeMap  map[string]string         // 表字段类型映射，入股为空则使用默认规则
	TplDir         string                    `validate:"required"` // 模板目录
	RootDir        string                    `validate:"required"` // 生成文件的根目录
	LayerDirMap    map[LayerName]string      // 各层级目录，如果为空则使用默认规则
	LayerNameMap   map[LayerName]LayerName   // 各层级名称，如果为空则使用默认规则
	LayerPrefixMap map[LayerName]LayerPrefix // 各层级前缀，如果为空则使用默认规则
}

type ControllerCfg struct {
	PackageName    string                    // 包名
	TargetFilename string                    // 目标文件名
	TplDir         string                    // 模板目录
	RootDir        string                    // 生成文件的根目录
	LayerDirMap    map[LayerName]string      // 各层级目录，如果为空则使用默认规则
	LayerNameMap   map[LayerName]LayerName   // 各层级名称，如果为空则使用默认规则
	LayerPrefixMap map[LayerName]LayerPrefix // 各层级前缀，如果为空则使用默认规则
}

type LayerName string

type LayerPrefix string

const (
	LayerNameRouter     LayerName = "router"
	LayerNameController LayerName = "controller"
	LayerNameService    LayerName = "service"
	LayerNameDto        LayerName = "dto"
	LayerNameRequest    LayerName = "request"
	LayerNameResponse   LayerName = "response"
	LayerNameModel      LayerName = "model"
	LayerNameErrorCode  LayerName = "errorCode"

	defaultLayerNameRequest  LayerName = "dto"
	defaultLayerNameResponse LayerName = "dto"

	defaultLayerPrefixController LayerPrefix = "ctr"
	defaultLayerPrefixService    LayerPrefix = "svc"
	defaultLayerPrefixDto        LayerPrefix = "dto"
	defaultLayerPrefixModel      LayerPrefix = "dao"
)

var defaultLayerPrefixMap = map[LayerName]LayerPrefix{
	LayerNameController: defaultLayerPrefixController,
	LayerNameService:    defaultLayerPrefixService,
	LayerNameModel:      defaultLayerPrefixModel,
	LayerNameDto:        defaultLayerPrefixDto,
}

var defaultLayerSpecialNameMap = map[LayerName]LayerName{
	LayerNameRequest:  defaultLayerNameRequest,
	LayerNameResponse: defaultLayerNameResponse,
}
