package glog

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

func getZapStandoutWriter() zapcore.WriteSyncer {
	return os.Stdout
}

func getZapFileWriter(cfg *LoggerConfig, fileSuffix string) (zapcore.WriteSyncer, error) {
	dir := strings.TrimSuffix(cfg.Dir, "/") + "/" + time.Now().Format("20060102")
	if ok := fileExists(dir); !ok {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
	logFilename := fmt.Sprintf("%s%s", cfg.service, fileSuffix)
	logFilepath := path.Join(dir, logFilename)
	file, openErr := os.OpenFile(logFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return nil, openErr
	}
	writer := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(file),
		Size:          256 * 1024,
		FlushInterval: time.Second * 5,
		Clock:         nil,
	}

	return writer, nil
}
