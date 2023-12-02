package log

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func init() {
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		Level:           log.InfoLevel,
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.StampMilli,
	})
}

func SetLevel(level string) error {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return err
	}
	Logger.SetLevel(lvl)
	return nil
}
