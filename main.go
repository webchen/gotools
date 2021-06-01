package main

import (
	"fmt"

	"github.com/webchen/gotools/base/conf"
)

func main() {
	f := conf.GetConfig("system.deployPathName", "path").(string)
	fmt.Println(f)
}
