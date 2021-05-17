/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2021-03-16 11:32:02
# File Name: messager.go
# Description:
####################################################################### */

// https://developers.dingtalk.com/document/app/custom-robot-access

package messager

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/ant-libs-go/http"
	http_cli "github.com/ant-libs-go/http/client"
	"github.com/ant-libs-go/util"
)

type MESSAGER string

const (
	MESSAGER_DINGDING_TEXT     MESSAGER = "dingding_text"
	MESSAGER_DINGDING_DOWNLOAD MESSAGER = "dingding_download"
)

var MESSAGER_DEFINED = map[MESSAGER]*struct {
	url        string
	restCodec  http.Codec
	restMethod http_cli.REST_METHOD
}{
	MESSAGER_DINGDING_TEXT: {
		url:        "https://oapi.dingtalk.com/robot/send?access_token=__TOKEN__",
		restCodec:  http.CODEC_JSON,
		restMethod: http_cli.REST_METHOD_JSON_POST},
	MESSAGER_DINGDING_DOWNLOAD: {
		url:        "https://oapi.dingtalk.com/robot/send?access_token=__TOKEN__",
		restCodec:  http.CODEC_JSON,
		restMethod: http_cli.REST_METHOD_JSON_POST},
}

type Messager struct {
	lock sync.RWMutex
	cfg  *Cfg
}

func NewMessager(cfg *Cfg) *Messager {
	o := &Messager{cfg: cfg}
	return o
}

func (this *Messager) Call(params map[string]string) (err error) {
	switch this.cfg.Messager {
	case MESSAGER_DINGDING_TEXT:
		err = this.callDingDingText(params)
	case MESSAGER_DINGDING_DOWNLOAD:
		err = this.callDingDingDown(params)
	default:
		err = fmt.Errorf("messager#%s not support", this.cfg.Messager)
	}
	return
}

func (this *Messager) callDingDingText(params map[string]string) (err error) {
	if len(params["text"]) > 18000 { // dingding max message length
		params["text"] = params["text"][:18000]
	}
	msg := `{"msgtype": "text", "text": {"content": "__TEXT__"}, "at": {"atMobiles": __ATS__}}`
	resp := &struct {
		Errcode int32
		Errmsg  string
	}{}
	_, err = this.buildRestClient().Call(nil, nil, this.buildMsg(msg, params), resp)
	if err == nil {
		util.IfDo(resp.Errcode != 0, func() { err = fmt.Errorf("reply code is %d, no 0. %s", resp.Errcode, resp.Errmsg) })
	}
	if err != nil {
		err = fmt.Errorf("push dingding text msg fail, %s", err)
	}
	return
}

func (this *Messager) callDingDingDown(params map[string]string) (err error) {
	msg := `{"msgtype": "actionCard", "actionCard": {
			"title": "__TITLE__",
			"text": "__TEXT__",
			"btnOrientation": "0",
			"btns": [{"title": "点击下载(截止:__EXPIRE_TIME__)", "actionURL": "__DOWNLOAD_URL__"}]
		}}`
	resp := &struct {
		Errcode int32
		Errmsg  string
	}{}
	_, err = this.buildRestClient().Call(nil, nil, this.buildMsg(msg, params), resp)
	if err == nil {
		util.IfDo(resp.Errcode != 0, func() { err = fmt.Errorf("reply code is %d, no 0. %s", resp.Errcode, resp.Errmsg) })
	}
	if err != nil {
		err = fmt.Errorf("push dingding download msg fail, %s", err)
	}
	return
}

func (this *Messager) buildRestClient() (r *http_cli.RestClientPool) {
	restCfg, _ := http_cli.LoadCfg(this.cfg.RestClient)
	if restCfg == nil {
		restCfg = &http_cli.Cfg{}
	}
	restCfg.Url = strings.Replace(MESSAGER_DEFINED[this.cfg.Messager].url, "__TOKEN__", this.cfg.Token, -1)
	restCfg.Codec = MESSAGER_DEFINED[this.cfg.Messager].restCodec
	restCfg.Method = MESSAGER_DEFINED[this.cfg.Messager].restMethod
	r = http_cli.NewRestClientPool(restCfg)
	return
}

func (this *Messager) buildMsg(msg string, params map[string]string) (r string) {
	macros := map[string]string{}
	b, _ := json.Marshal(this.cfg.Ats)
	macros["__ATS__"] = string(b)
	for k, v := range params {
		macros[fmt.Sprintf("__%s__", strings.ToUpper(k))] = v
	}

	t := []byte(msg)
	for k, v := range macros {
		t = util.BytesReplace(t, []byte(k), []byte(v), -1)
	}
	r = string(t)
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
