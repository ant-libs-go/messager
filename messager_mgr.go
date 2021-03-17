/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2021-03-15 21:00:53
# File Name: messager_mgr.go
# Description:
####################################################################### */

package messager

import (
	"fmt"
	"sync"

	"github.com/ant-libs-go/config"
	"github.com/ant-libs-go/config/options"
)

var (
	once  sync.Once
	lock  sync.RWMutex
	pools map[string]*Messager
)

type messagerConfig struct {
	Cfgs map[string]*Cfg `toml:"messager"`
}

type Cfg struct {
	Messager   MESSAGER `toml:"messager"`
	RestClient string   `toml:rest_client` // 指定rest配置，非必须
	Token      string   `toml:"token"`
	Ats        []string `toml:"ats"`
}

func Call(name string, params map[string]string) (err error) {
	var cli *Messager
	cli, err = getCli(name)
	if err == nil {
		err = cli.Call(params)
	}
	return
}

func getCli(name string) (r *Messager, err error) {
	lock.RLock()
	r = pools[name]
	lock.RUnlock()
	if r == nil {
		r, err = addCli(name)
	}
	return
}

func addCli(name string) (r *Messager, err error) {
	var cfg *Cfg
	if cfg, err = LoadCfg(name); err != nil {
		return
	}
	r = NewMessager(cfg)

	lock.Lock()
	pools[name] = r
	lock.Unlock()
	return
}

func LoadCfg(name string) (r *Cfg, err error) {
	var cfgs map[string]*Cfg
	if cfgs, err = loadCfgs(); err != nil {
		return
	}
	if r = cfgs[name]; r == nil {
		err = fmt.Errorf("messager#%s not configed", name)
		return
	}
	return
}

func loadCfgs() (r map[string]*Cfg, err error) {
	r = map[string]*Cfg{}

	cfg := &messagerConfig{}
	once.Do(func() {
		_, err = config.Load(cfg, options.WithOnChangeFn(func(cfg interface{}) {
			lock.Lock()
			defer lock.Unlock()
			pools = map[string]*Messager{}
		}))
	})

	cfg = config.Get(cfg).(*messagerConfig)
	if err == nil && (cfg.Cfgs == nil || len(cfg.Cfgs) == 0) {
		err = fmt.Errorf("not configed")
	}
	if err != nil {
		err = fmt.Errorf("messager load cfgs error, %s", err)
		return
	}
	r = cfg.Cfgs
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
