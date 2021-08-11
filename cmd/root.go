package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zcx2001/goshare/logger"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "webDownload",
	Short: "webDownload",
	Long:  "webDownload",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		viper.AddConfigPath("conf/")
		viper.SetConfigName("app")

		// 把环境变量内的_替换成. 并自动导入
		viper.SetEnvPrefix("APP")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}

		// 初始化日志框架
		logger.Init()

		// 初始化 maxprocs
		_, _ = maxprocs.Set(maxprocs.Logger(func(s string, i ...interface{}) {
			logger.Log.Debug("maxprocs logger", zap.String("info", fmt.Sprintf(s, i...)))
		}))
	})
}
