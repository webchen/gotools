package main

import (
	"fmt"
	"time"

	"github.com/webchen/gotools/base/conf"
)

func main() {
	t, _ := time.ParseDuration(conf.GetConfig("system.http.queryTimeOut", "3s").(string))
	fmt.Printf("%+v", t)
}
