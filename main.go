package main

import (
	"bytes"
	"os"

	"github.com/jexia-com/logtest/adapter"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		buffer   = adapter.NewBuffer(0)
		writer   bytes.Buffer
		failing  = adapter.NewFailingWriter(&writer, 0)
		buffered = adapter.NewBufferingWriter(failing, buffer)
		logger   = &logrus.Logger{
			Out:       buffered,
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
			ExitFunc: func(code int) {
				buffer.Flush(os.Stdout)

				os.Exit(code)
			},
			ReportCaller: false,
		}
	)

	// defer buffer.Flush(os.Stdout)

	logger.Debug("DEBUG message")
	logger.Info("INFO message")
	logger.Warning("WARNING message")
	logger.Error("ERROR message")
	logger.Fatal("FATAL message")
}
