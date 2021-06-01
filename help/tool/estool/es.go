package estool

import (
	"bytes"
	"context"
	"log"
	"strings"
	"time"

	"github.com/webchen/gotools/base/conf"
	"github.com/webchen/gotools/base/jsontool"
	"github.com/webchen/gotools/help/tool/nettool"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var es *elasticsearch.Client

func init() {
	var hostEmpty []interface{}
	host := conf.GetConfig("es.host", hostEmpty).([]interface{})
	var hostList []string
	for _, v := range host {
		hostList = append(hostList, v.(string))
	}
	user := conf.GetConfig("es.user", "").(string)
	password := conf.GetConfig("es.password", "").(string)
	cfg := elasticsearch.Config{
		Addresses: hostList,
		Username:  user,
		Password:  password,
		// ...
	}
	var err error
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("无法初始化ES类  [%+v]", err)
	}
}

// WriteLog 往ES里面写LOG
func WriteLog(level string, message string, v ...interface{}) {
	index := (conf.GetConfig("es.index", "gateway_pub")).(string)
	go (func() {
		data := map[string]interface{}{
			"@timestamp": time.Now().Format(time.RFC3339Nano),
			"level":      level,
			"ip":         nettool.GetLocalFirstIPStr(),
			"message":    message,
			"content":    v,
		}
		body := jsontool.MarshalToString(data)
		req := esapi.IndexRequest{
			Index:   index,
			Body:    bytes.NewReader([]byte(body)),
			Refresh: "true",
		}
		res, err := req.Do(context.Background(), es)
		if err != nil || res == nil {
			log.SetPrefix("ESERROR")
			log.Printf("write log error [%+v] [%+v]", data, err)
			return
		}
		defer res.Body.Close()
		if strings.Contains(res.String(), "error") {
			log.Println(res.String())
		}
	})()
}
