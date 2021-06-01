package conf

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/webchen/gotools/base"
	"github.com/webchen/gotools/base/dirtool"
	"github.com/webchen/gotools/base/jsontool"

	"github.com/zouyx/agollo/v4"
	apolloConfig "github.com/zouyx/agollo/v4/env/config"
)

// 全局配置变量
var config = make(map[string]map[string]interface{})

var apolloData map[string]string

func init() {
	loadApolloConfig()
	initApollo()
	initLocal()

	//os.Exit(1)
}

func initLocal() {
	dir := dirtool.GetConfigPath()
	fileList, _ := filepath.Glob(filepath.Join(dir, "*"))
	for j := range fileList {
		ext := path.Ext(fileList[j])
		if ext == ".json" {
			fileName := strings.ReplaceAll(strings.ReplaceAll(fileList[j], filepath.Dir(fileList[j])+string(os.PathSeparator), ""), ext, "")
			conf := make(map[string]interface{})
			jsontool.LoadFromFile(fileList[j], &conf)
			config[fileName] = conf
		}
	}
}

func loadApolloConfig() {
	f := dirtool.GetConfigPath() + "apollo.json"
	exists, _ := dirtool.PathExist(f)
	if !exists {
		return
	}
	jsontool.LoadFromFile(f, &apolloData)
}

func initApollo() {
	if apolloData == nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			base.LogPanic("initApollo error", p)
		}
	}()

	c := &apolloConfig.AppConfig{
		AppID:         apolloData["appID"],
		Cluster:       apolloData["cluster"],
		IP:            apolloData["host"],
		NamespaceName: apolloData["namespace"],
		Secret:        apolloData["secret"],
	}
	//	agollo.SetLogger(&log.DefaultLogger{})
	client, _ := agollo.StartWithConfig(func() (*apolloConfig.AppConfig, error) {
		return c, nil
	})

	cache := client.GetConfigCache(c.NamespaceName)
	cache.Range(func(key, value interface{}) bool {
		configFilePath := dirtool.GetConfigPath() + key.(string) + ".json"
		ioutil.WriteFile(configFilePath, []byte(value.(string)), 0777)
		//fmt.Printf("key: %+v   val:%+v\n", key, value)
		return true
	})

	//	value, _ := cache.Get("es")
	//	fmt.Printf("%+v\n%+v\n", cache, value)
}

// GetConfig 获取JSON的配置，key支持"."操作，如：GetConfig("conf.runtime")，表示获取conf.json文件里面，runtime的值
func GetConfig(key string, def interface{}) interface{} {
	defer func() {
		recover()
	}()
	arr := strings.Split(key, ".")
	if len(arr) == 0 {
		return def
	}
	if len(arr) == 1 {
		if config[arr[0]] == nil {
			return def
		}
		return config[arr[0]]
	}
	confDeep := config[arr[0]][arr[1]]
	if len(arr) == 2 {
		if confDeep == nil {
			return def
		}
		return confDeep
	}
	for j := 2; j < len(arr); j++ {
		c, _ := confDeep.(interface{})
		if c == nil {
			return def
		}
		confDeep = confDeep.(map[string]interface{})[arr[j]]
		if confDeep == nil {
			return def
		}
	}
	return confDeep
}
