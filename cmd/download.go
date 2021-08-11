package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/zcx2001/goshare/httpService"
	"github.com/zcx2001/webDownload/pkg/service"
	"net/http"
	"os"
	"os/signal"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "download [flags] weburl",
	Short: "download web",
	Long:  "download web",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			err = errors.New("args size error")
			return
		}
		return
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)

		err = service.DownloadHtml(args[0], true)

		// 初始化网站服务
		web := httpService.New(
			httpService.WithPort("8080"),
			httpService.WithStatic("/", http.Dir("demo.mxyhn.xyz:8020")),
		)
		web.Start()

		<-quit

		// 结束网站服务
		web.Stop()

		return
	},
}
