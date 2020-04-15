package robust

import "github.com/golang/glog"

//处理异常等级
const (
	Info 	= iota
	Warning
	Error 
	Fatal
	Exit
)

//Recover 捕捉异常
func Recover(level int, arr ...ErrorHandle) {

	if r := recover(); r != nil {
		switch level {
		case Info:		glog.Info(r)
		case Warning:	glog.Warning(r)
		case Error:		glog.Error(r)
		case Fatal:		glog.Fatal(r)
		case Exit:		glog.Exit(r)
		}
		for _, handle := range arr {
			handle(r)
		}
	}
}

//ErrorHandle 捕捉到异常后，处理
type ErrorHandle func(interface{})