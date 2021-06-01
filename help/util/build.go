package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/webchen/gotools/base"
	"github.com/webchen/gotools/base/dirtool"
	"github.com/webchen/gotools/help/tool/conf"
)

//var buildDir = "build" + string(os.PathSeparator)
var buildDir = ""
var osList = [4]string{"linux", "windows", "mac", "freebsd"}
var fileName = "carGatewayService"

// DoBuild 构建
func DoBuild(osName string) {
	has := false
	for _, v := range osList {
		if osName == v {
			has = true
			break
		}
	}

	if !has {
		base.LogPanic("system : "+osName+" was not supported ...", nil)
		return
	}

	dir := dirtool.GetBasePath()
	fileList, _ := filepath.Glob(filepath.Join(dir, "*"))
	includeFile := ""
	_, file, _, _ := runtime.Caller(0)
	for j := range fileList {
		str := strings.ReplaceAll(fileList[j], "\\", "/")
		if !strings.EqualFold(str, file) {
			if str[len(str)-3:] == ".go" {
				includeFile += str + " "
			}
		}
	}
	deployConf := conf.GetConfig("system.deployPathName", "public").(string)
	if deployConf == "" {
		deployConf = "public"
	}
	deployPath := dirtool.GetParentDirectory(dirtool.GetParentDirectory(dir)) + string(os.PathSeparator) + deployConf + string(os.PathSeparator) + buildDir
	dirtool.MustCreateDir(deployPath)
	deployConfigPath := deployPath + "config" + string(os.PathSeparator)
	dirtool.MustCreateDir(deployConfigPath)

	confPath := dirtool.GetConfigPath()
	confList, _ := ioutil.ReadDir(confPath)
	for _, f := range confList {
		fsBytes, _ := ioutil.ReadFile(confPath + f.Name())
		info := string(fsBytes)
		err := ioutil.WriteFile(deployConfigPath+f.Name(), []byte(info), 0777)
		if err != nil {
			panic(err)
		}
	}

	if osName == "windows" {
		fileName += ".exe"
	}

	sys := runtime.GOOS
	tmpFile := deployPath + "tmp"
	if sys == "windows" {
		tmpFile += ".bat"
	}
	if sys == "linux" {
		tmpFile += ".sh"
	}
	cmdStr := getCmd(osName, deployPath+fileName, includeFile)
	ioutil.WriteFile(tmpFile, []byte(cmdStr), 0666)

	cmd := exec.Command(tmpFile)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("build error : %+v\n", err)
	} else {
		fmt.Printf("build success ...\ndirectory : %s\n", deployPath)
	}
	os.Remove(tmpFile)

}

func getCmd(osName string, fileName string, files string) string {

	cmd := fmt.Sprintf(
		`
SET CGO_ENABLED=0
SET GOOS=%s
SET GOARCH=amd64
go build -i -o %s %s
`, osName, fileName, files)
	return cmd
}
