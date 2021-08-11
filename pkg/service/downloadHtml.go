package service

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zcx2001/goshare/logger"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

func DownloadHtml(urlAddr string, rewrite bool) (err error) {
	u, err := url.ParseRequestURI(urlAddr)
	if err != nil {
		return
	}

	bs, err := downloadFile(u, rewrite)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bs))
	if err != nil {
		return
	}

	findLink(u, doc)

	findImg(u, doc)

	findScript(u, doc)

	findA(u, doc)

	return
}

func findA(addr *url.URL, doc *goquery.Document) {
	// 找出全部的 a
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if value, exists := s.Attr("href"); exists {
			if strings.HasPrefix(value, "http") {
				logger.Log.Debug("Discover the outbound links", zap.String("url", value))
			} else if strings.HasPrefix(value, "#") {

			} else {
				imgurl, err := url.Parse(value)
				if err != nil {
					logger.Log.Error("img src error", zap.String("src", value))
					return
				}
				logger.Log.Debug("Discover A", zap.String("url", addr.ResolveReference(imgurl).String()))
				_ = DownloadHtml(addr.ResolveReference(imgurl).String(), false)
			}
		}
	})
}

func findLink(addr *url.URL, doc *goquery.Document) {
	// 找出全部的 script
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		if value, exists := s.Attr("href"); exists {
			if strings.HasPrefix(value, "http") {
				logger.Log.Debug("Discover the outbound links", zap.String("href", value))
			} else {
				linkurl, err := url.Parse(value)
				if err != nil {
					logger.Log.Error("link src error", zap.String("href", value))
					return
				}
				logger.Log.Debug("Discover link", zap.String("href", addr.ResolveReference(linkurl).String()))
				cssData, _ := downloadFile(addr.ResolveReference(linkurl), false)

				if path.Ext(value) == ".css" {
					//fmt.Println("css =", value)
					findCSSurl(addr.ResolveReference(linkurl), cssData)
				}
			}
		}
	})
}

func findCSSurl(addr *url.URL, cssData []byte) {
	re := regexp.MustCompile("url\\((\\S+?)\\)")
	data := re.FindAllSubmatch(cssData, -1)
	for _, d1 := range data {
		for i, d2 := range d1 {
			if i != 0 {
				value := strings.Trim(string(d2), "\"'")

				if strings.HasPrefix(value, "http") {
					logger.Log.Debug("Discover the outbound links", zap.String("href", value))
				} else {
					linkurl, err := url.Parse(value)
					if err != nil {
						logger.Log.Error("link src error", zap.String("href", value))
						return
					}
					logger.Log.Debug("Discover link", zap.String("href", addr.ResolveReference(linkurl).String()))
					_, _ = downloadFile(addr.ResolveReference(linkurl), false)
				}
			}
		}
	}
}

func findScript(addr *url.URL, doc *goquery.Document) {
	// 找出全部的 script
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if value, exists := s.Attr("src"); exists {
			if strings.HasPrefix(value, "http") {
				logger.Log.Debug("Discover the outbound links", zap.String("url", value))
			} else {
				scripturl, err := url.Parse(value)
				if err != nil {
					logger.Log.Error("script src error", zap.String("src", value))
					return
				}
				logger.Log.Debug("Discover script", zap.String("url", addr.ResolveReference(scripturl).String()))
				_, _ = downloadFile(addr.ResolveReference(scripturl), false)
			}
		}
	})
}

func findImg(addr *url.URL, doc *goquery.Document) {
	// 找出全部的 img
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		if value, exists := s.Attr("src"); exists {
			if strings.HasPrefix(value, "http") {
				logger.Log.Debug("Discover the outbound links", zap.String("url", value))
			} else if strings.HasPrefix(value, "data:") {
				logger.Log.Debug("Discover base64 data", zap.String("data", value))
			} else {
				imgurl, err := url.Parse(value)
				if err != nil {
					logger.Log.Error("img src error", zap.String("src", value))
					return
				}
				logger.Log.Debug("Discover image", zap.String("url", addr.ResolveReference(imgurl).String()))
				_, _ = downloadFile(addr.ResolveReference(imgurl), false)
			}
		}
	})

}

func downloadFile(addr *url.URL, rewrite bool) (bs []byte, err error) {
	filename := path.Join(addr.Host, path.Dir(addr.Path), path.Base(addr.Path))

	err = os.MkdirAll(path.Dir(filename), os.ModePerm)
	if err != nil {
		logger.Log.Error("mkdirAll error", zap.Error(err))
		return
	}

	if !rewrite {
		if _, err = os.Stat(filename); !os.IsNotExist(err) {
			logger.Log.Debug("file is exist", zap.String("filename", filename))
			return
		}
	}

	resp, err := http.Get(addr.String())
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != 200 {
		err = fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
		return
	}

	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("ioutil ReadAll", zap.Error(err))
		return
	}

	err = ioutil.WriteFile(filename, bs, os.ModePerm)
	if err != nil {
		logger.Log.Error("ioutil WriteFile", zap.Error(err))
		return
	}

	return
}
