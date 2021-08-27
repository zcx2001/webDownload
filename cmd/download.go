package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/zcx2001/goshare/httpService"
	"github.com/zcx2001/goshare/logger"
	"github.com/zcx2001/webDownload/pkg/service"
	"net"
	"net/http"
	"net/url"
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

		urlAddr, err := url.ParseRequestURI(args[0])
		if err != nil {
			return
		}

		err = service.DownloadHtml(urlAddr, true)

		logger.Log.Debug("download ok")

		host, _, err := net.SplitHostPort(urlAddr.Host)
		if err != nil {
			return
		}
		// 初始化网站服务
		web := httpService.New(
			httpService.WithPort("8080"),
			httpService.WithStatic("/", http.Dir(host)),
		)
		web.Start()

		<-quit

		// 结束网站服务
		web.Stop()

		return
	},
}
