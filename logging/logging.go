package logging

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
	"runtime"
)

type Logger struct {
	*logrus.Logger
}

func GetLogger(ctx context.Context) *Logger {
	return LoggerFromContext(ctx)
}

func NewLogger() *Logger {
	l := logrus.New()

	l.SetReportCaller(true)
	setTextFormat(l)

	l.SetLevel(logrus.InfoLevel)
	setLumberjackOutput(l)

	return &Logger{
		l,
	}
}

func setTextFormat(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors: true,
		FullTimestamp: true,
	})
}

func setLumberjackOutput(logger *logrus.Logger) {
	logger.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/app.log", // Основной файл логов
		MaxSize:    3,                // Макс. размер файла (МБ)
		MaxBackups: 3,                // Количество старых файлов
		MaxAge:     30,               // Дни хранения
		Compress:   true,             // Сжатие старых логов
	})
}
