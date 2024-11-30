package generate

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/morehao/go-tools/conf"
	"github.com/morehao/go-tools/dbClient"
	"github.com/morehao/go-tools/glog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:embed template
var templatesFS embed.FS

var workDir string
var cfg *Config
var MysqlClient *gorm.DB

// Cmd represents the generate command
var Cmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code based on templates",
	Long:  `Generate code for different layers like module, model, and API based on predefined templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化配置和 MySQL 客户端
		if cfg == nil {
			workDir, _ := os.Getwd()
			configFilepath := filepath.Join(workDir, "config", "config.yaml")
			conf.LoadConfig(configFilepath, &cfg)
		}

		// 初始化日志
		if err := glog.NewLogger(&cfg.Log, glog.WithZapOptions(zap.AddCallerSkip(3))); err != nil {
			panic("glog initZapLogger error")
		}

		// 延迟初始化 Mysql 客户端
		if MysqlClient == nil {
			mysqlClient, getMysqlClientErr := dbClient.InitMysql(cfg.Mysql)
			if getMysqlClientErr != nil {
				panic("get mysql client error")
			}
			MysqlClient = mysqlClient
		}

		mode, _ := cmd.Flags().GetString("mode")
		workDir, _ := os.Getwd()

		if workDir == "" {
			fmt.Println("Please provide a working directory using --workdir flag")
			return
		}

		switch mode {
		case "module":
			if err := genModule(); err != nil {
				fmt.Printf("Error generating module: %v\n", err)
				return
			}
			fmt.Println("Module generated successfully")
		case "model":
			if err := genModel(); err != nil {
				fmt.Printf("Error generating model: %v\n", err)
				return
			}
			fmt.Println("Model generated successfully")
		case "api":
			if err := genApi(); err != nil {
				fmt.Printf("Error generating api: %v\n", err)
				return
			}

		// 这里可以添加其他模式的处理逻辑
		default:
			fmt.Println("Invalid mode. Available modes are: module, model, api")
		}
	},
}

func init() {
	// 定义 generate 命令的参数
	Cmd.Flags().StringP("mode", "m", "", "Mode of code generation (module, model, api)")
}
