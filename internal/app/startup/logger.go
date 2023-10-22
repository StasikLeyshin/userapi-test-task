package startup

import (
	"github.com/sirupsen/logrus"
	logger "refactoring/pkg"
)

func NewLogger() *logrus.Logger {
	return logger.NewLogger()
}
