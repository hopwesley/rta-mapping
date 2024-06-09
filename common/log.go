package common

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
)

const (
	defaultLogFileName = "rta.log"
)

var _logInstance *logrus.Logger
var logOnce sync.Once

func LogInst() *logrus.Logger {
	logOnce.Do(func() {

		log := logrus.New()

		log.SetOutput(&lumberjack.Logger{
			Filename:   defaultLogFileName, // 日志文件路径
			MaxSize:    30,                 // 文件最大 MB
			MaxBackups: 10,                 // 最大备份文件个数
			MaxAge:     28,                 // 文件最大保存天数
			Compress:   true,               // 是否压缩/归档旧文件
		})
		_logInstance = log
	})
	return _logInstance
}
