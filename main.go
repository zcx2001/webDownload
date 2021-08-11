package main

import (
	"github.com/zcx2001/webDownload/cmd"
	"math/rand"

	"time"
)

func main() {
	//初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	cmd.Execute()
}
