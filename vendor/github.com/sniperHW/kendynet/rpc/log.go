package rpc

import(
	"github.com/sniperHW/kendynet/golog"
	"sync/atomic"	
	"fmt"
)

var logger *golog.Logger
var is_init int32

func InitLogger(out *golog.OutputLogger,name ...string) {
	if atomic.CompareAndSwapInt32(&is_init, 0, 1) {
		var fullname string
		if len(name) > 0 {
			fullname = fmt.Sprintf("(%s|rpc)",name[0])
		} else {
			fullname = "(rpc)"
		}
		logger = golog.New(fullname,out)
		logger.Debugf("%s logger init",fullname)
	}	
}

func Debugf(format string, v ...interface{}) {
	if nil != logger {
		logger.Debugf(format,v...)
	}
}

func Debugln(v ...interface{}) {
	if nil != logger {
		logger.Debugln(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if nil != logger {
		logger.Infof(format,v...)
	}
}

func Infoln(v ...interface{}) {
	if nil != logger {
		logger.Infoln(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if nil != logger {
		logger.Warnf(format,v...)
	}
}

func Warnln(v ...interface{}) {
	if nil != logger {
		logger.Warnln(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if nil != logger {
		logger.Errorf(format,v...)
	}
}

func Errorln(v ...interface{}) {
	if nil != logger {
		logger.Errorln(v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	if nil != logger {
		logger.Fatalf(format,v...)
	}
}

func Fatalln(v ...interface{}) {
	if nil != logger {
		logger.Fatalln(v...)
	}
}