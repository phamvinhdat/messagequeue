package messagequeue

import (
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006/01/02 15:04:05.000",
		DisableLevelTruncation: true,
	})
}

func WithDebug() {
	logrus.SetLevel(logrus.DebugLevel)
}
