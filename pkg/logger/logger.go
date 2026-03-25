package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() {

	Log.SetOutput(os.Stdout)

	Log.SetFormatter(&logrus.JSONFormatter{})

	Log.SetLevel(logrus.InfoLevel)
}
