package logs

import (
	"fmt"
	"strings"

	base "github.com/webchen/gotools/base"
	"github.com/webchen/gotools/base/conf"
	"github.com/webchen/gotools/base/jsontool"
	"github.com/webchen/gotools/enum/def"

	"log"
	"os"
	"time"
	//	log "github.com/sirupsen/logrus"
)

//var maxStack = 20
//var separator = "---------------------------------------" + fmt.Sprintln()

// debug -> 0 info/readmq -> 1 warning/query -> 2 error/message -> 3 critial -> 4 exit -> 9

var fileLogger *log.Logger
var cmdLogger *log.Logger

// 日志等级
var cmdLevel float64 = 0
var fileLevel float64 = 0

//var esLevel float64 = 0

func init() {
	cmdLevel = (conf.GetConfig("conf.logCmdLevel", 0.0)).(float64)
	fileLevel = (conf.GetConfig("conf.logFileLevel", 0.0)).(float64)
	//esLevel = (conf.GetConfig("conf.esLevel", 0.0)).(float64)

	fileLogger = access("log")

	cmdLogger = newCmdLogger("")

}

// 初始化cmd环境下的logger对象
func newCmdLogger(level string) *log.Logger {
	l := new(log.Logger)
	l.SetPrefix("[" + level + "] ")
	l.SetFlags(log.Lmicroseconds)

	l.SetOutput(os.Stdout)
	return l
}

// Debug 日志
func Debug(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)

	fileLogger.SetPrefix("[Debug] ")
	cmdLogger.SetPrefix("[Debug] ")

	if fileLevel == 0 {
		fileLogger.Println(s)
	}

	if cmdLevel == 0 {
		cmdLogger.Println(s)
	}
}

// Info 日志
func Info(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	fileLogger.SetPrefix("[info] ")
	cmdLogger.SetPrefix("[info] ")

	if fileLevel <= 1 {
		fileLogger.Println(s)
	}

	if cmdLevel <= 1 {
		cmdLogger.Println(s)
	}
}

// Warning 日志
func Warning(message string, data interface{}, withTrace bool) {
	info := &base.LogObj{}
	info.Message = message
	info.Time = time.Now().Format(def.FullTimeMicroFormat)
	info.Level = "Warning"
	if withTrace {
		info.Trace = Trace(message)
	} else {
		info.Trace = ""
	}

	s := jsontool.MarshalToString(info)
	if fileLevel <= 2 {
		fileLogger.SetPrefix("[Warning] ")
		fileLogger.SetFlags(0)
		fileLogger.Println(s)
	}

	if cmdLevel <= 2 {
		cmdLogger.SetPrefix("[Warning] ")
		cmdLogger.SetFlags(0)
		cmdLogger.Println(s)
	}
}

// Error 日志
func Error(message string, data interface{}) {
	info := &base.LogObj{}
	info.Message = message
	info.Time = time.Now().Format(def.FullTimeMicroFormat)
	info.Level = "Error"
	info.Trace = Trace(message)

	s := jsontool.MarshalToString(info)
	if fileLevel <= 3 {
		fileLogger.SetPrefix("[Error] ")
		fileLogger.SetFlags(0)
		fileLogger.Println(s)
	}

	if cmdLevel <= 3 {
		cmdLogger.SetPrefix("[Error] ")
		cmdLogger.Println(s)
	}
}

// Message 日志
func Message(message string, data interface{}, withTrace bool) {
	info := &base.LogObj{}
	info.Message = message
	info.Time = time.Now().Format(def.FullTimeMicroFormat)
	info.Level = "Message"
	if withTrace {
		info.Trace = Trace(message)
	} else {
		info.Trace = ""
	}

	s := jsontool.MarshalToString(info)

	if fileLevel <= 4 {
		fileLogger.SetPrefix("[Message] ")
		fileLogger.SetFlags(0)
		fileLogger.Println(s)
	}

	if cmdLevel <= 4 {
		cmdLogger.SetPrefix("[Message] ")
		cmdLogger.Println(s)
	}
}

// MessageClient 日志
func MessageClient(message string, data interface{}, withTrace bool) {

	info := &base.LogObj{}
	info.Message = message
	info.Time = time.Now().Format(def.FullTimeMicroFormat)
	info.Level = "MessageClient"
	if withTrace {
		info.Trace = Trace(message)
	} else {
		info.Trace = ""
	}
	val, isString := data.(string)
	if isString && strings.TrimSpace(val) != "" {
		info.Data = jsontool.JSONStrFormat(val)
	} else {
		if err, isErr := data.(error); isErr {
			info.Data = err.Error()
		} else {
			info.Data = jsontool.JSONStrFormat(jsontool.MarshalToString(data))
		}
	}

	s := jsontool.MarshalToString(info)

	if fileLevel <= 4 {
		fileLogger.SetPrefix("[MessageClient] ")
		fileLogger.SetFlags(0)
		fileLogger.Println(s)
	}

	if cmdLevel <= 4 {
		cmdLogger.SetPrefix("[MessageClient] ")
		cmdLogger.Println(s)
	}
}

// Critial 日志
func Critial(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if fileLevel <= 4 {
		fileLogger.SetPrefix("[CRITIAL] ")
		fileLogger.Println(s)
	}

	if cmdLevel <= 4 {
		cmdLogger.SetPrefix("[Critial] ")
		cmdLogger.Println(s)
	}
	//emailtool.SendAlertEmail(s)
	//runtime.Goexit()
}

// Exited 手动EXIT的时候打印日志
func Exited(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if fileLevel <= 9 {
		fileLogger.SetPrefix("[EXIT] ")
		fileLogger.Println(s)
	}

	if cmdLevel <= 9 {
		cmdLogger.Println(s)
	}
	//emailtool.SendAlertEmail(s)
	os.Exit(-1)
	//	runtime.Goexit()
}

// ExitedProcess 如果err!=nil，则打印日志，并退出
func ExitedProcess(err error, msg string) {
	if err != nil {
		msg += "\n%s"
		Exited(msg, Trace(err))
		panic(err)
	}
}

// ErrorProcess 错误处理
func ErrorProcess(err error, msg string) bool {
	if err != nil {
		msg += "\n"
		Error(msg, err)
		return true
	}
	return false
}

// CritialProcess 错误处理
func CritialProcess(err error, msg string) bool {
	if err != nil {
		msg += "\n%s"
		Critial(msg, Trace(err))
		return true
	}
	return false
}

// Query 请求第三方的日志
func Query(format string, v ...interface{}) {
	if fileLevel <= 2 {
		fileLogger.SetPrefix("[QUERY] ")
		fileLogger.Println(fmt.Sprintf(format, v...))
	}
}

// Show 打印一定会显示的信息（用于系统层面）
func Show(format string, v ...interface{}) {
	cmdLogger.SetPrefix("[show] ")
	cmdLogger.Println(fmt.Sprintf(format, v...))
	//log.Info(fmt.Sprintf(format, v...))
}

// ES 只打在ES里面的日志
/*
func ES(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	estool.WriteLog("es", s, traceInfo(v...))
}
*/

func access(fileName string) (l *log.Logger) {
	return base.CreateLogFileAccess(fileName)
}

// Trace 对外TRACE
func Trace(v ...interface{}) string {
	return base.TraceInfo(v)
}

// WebAccess WEB端访问日志
func WebAccess(format string, v ...interface{}) {
	open := conf.GetConfig("conf.openWebAccessLog", false).(bool)
	if !open {
		return
	}
	data := fmt.Sprintf(format, v...)
	if strings.Contains(data, "kube-probe/") || strings.Contains(data, "SLBHealthCheck") {
		return
	}
	fileLogger.SetPrefix("[ACCESS] ")
	fileLogger.Println(data)
}
