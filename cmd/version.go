package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zcx2001/goshare/logger"
	"github.com/zcx2001/webDownload/pkg/version"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Long:  "version",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.Info("app", zap.String("version", version.VERSION))
	},
}
