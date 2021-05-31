package base

import (
	"flag"
	"strings"
)

var builded = false
var buildOs = ""

// IsBuild 是否编译
func IsBuild() bool {
	return strings.TrimSpace(BuildOsName()) != ""
}

// BuildOsName 要编译的系统名称
func BuildOsName() string {
	if builded {
		return buildOs
	}
	flag.StringVar(&buildOs, "buildOs", "", "1) linux (default) \n 2) windows \n 3) mac \n 4) freebsd")
	flag.Parse()
	builded = true
	return buildOs
}
