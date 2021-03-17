/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2021-03-16 19:39:13
# File Name: messager/messager_test.go
# Description:
####################################################################### */

package messager

import (
	"fmt"
	"os"
	"testing"

	"github.com/ant-libs-go/config"
	"github.com/ant-libs-go/config/options"
	"github.com/ant-libs-go/config/parser"
	//. "github.com/smartystreets/goconvey/convey"
)

var globalCfg *config.Config

func TestMain(m *testing.M) {
	config.New(parser.NewTomlParser(),
		options.WithCfgSource("./test.toml"),
		options.WithCheckInterval(1))
	os.Exit(m.Run())
}

func TestDingDingText(t *testing.T) {
	fmt.Println(Call("dingdingtext", map[string]string{
		"text": "测试测试，文字消息测试",
	}))
}

func TestDingDingDownload(t *testing.T) {
	fmt.Println(Call("dingdingdownload", map[string]string{
		"title":        "测试标题",
		"text":         "测试测试，开始下载",
		"download_url": "http://www.baidu.com",
	}))
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
