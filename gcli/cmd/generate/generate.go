package generate

import (
	"embed"
	"fmt"
	"github.com/morehao/go-tools/conf"
	"github.com/morehao/go-tools/dbClient"
	"github.com/morehao/go-tools/glog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

//go:embed templates/*
var templatesFS embed.FS

// GenerateCmd represents the generate command
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code based on templates",
	Long:  `Generate code for different layers like module, model, and API based on predefined templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		mode, _ := cmd.Flags().GetString("mode")
		workDir, _ := os.Getwd()
		if workDir == "" {
			fmt.Println("Please provide a working directory using --workdir flag")
			return
		}

		switch mode {
		case "module":
			if err := genModule(workDir); err != nil {
				fmt.Printf("Error generating module: %v\n", err)
				return
			}
			fmt.Println("Module generated successfully")
		case "model":
			if err := genModel(workDir); err != nil {
				fmt.Printf("Error generating model: %v\n", err)
				return
			}
			fmt.Println("Model generated successfully")
		// 这里可以添加其他模式的处理逻辑
		default:
			fmt.Println("Invalid mode. Available modes are: module, model, api")
		}
	},
}

var Cfg *Config
var MysqlClient *gorm.DB

func init() {
	// 定义 generate 命令的参数
	GenerateCmd.Flags().StringP("mode", "m", "", "Mode of code generation (module, model, api)")

	// 初始化配置
	if workDir, err := os.Getwd(); err != nil {
		panic("get work dir error")
	} else {
		conf.SetAppRootDir(filepath.Join(workDir, "/internal/genCode"))
	}
	configFilepath := conf.GetAppRootDir() + "/config/config.yaml"

	conf.LoadConfig(configFilepath, &Cfg)

	// 初始化日志组件
	if err := glog.NewLogger(&Cfg.Log, glog.WithZapOptions(zap.AddCallerSkip(3))); err != nil {
		panic("glog initZapLogger error")
	}
	mysqlClient, getMysqlClientErr := dbClient.InitMysql(Cfg.Mysql)
	if getMysqlClientErr != nil {
		panic("get mysql client error")
	}
	MysqlClient = mysqlClient
}
