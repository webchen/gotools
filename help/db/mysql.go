package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/webchen/gotools/help/logs"
	"github.com/webchen/gotools/help/tool/conf"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlList = make(map[string]*sql.DB)

func init() {
	list := make(map[string]interface{})
	list = (conf.GetConfig("mysql", list)).(map[string]interface{})
	for k, v := range list {
		vv := v.(map[string]string)
		dsn := vv["connectionString"] //(conf.GetConfig("mysql."+k+".connectionString", "")).(string)
		obj, err := sql.Open("mysql", dsn)

		logs.ExitedProcess(err, fmt.Sprintf("无法连接[%s]数据库", k))

		obj.SetMaxOpenConns(6)      // 设置数据库的最大连接数
		obj.SetMaxIdleConns(6)      // 设置数据库的最大空闲连接数
		obj.SetConnMaxLifetime(120) //连接最长存活期，超过这个时间连接将不再被复用

		mysqlList[k] = obj
	}
}

// Get 获取MYSQL链接对象
func Get(name string) *sql.DB {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	if obj, exists := mysqlList[name]; exists {
		return obj
	}
	return nil
}
