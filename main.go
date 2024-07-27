package main

import (
	"myapp/api"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02/01/2006.15:04:05",
	})
	logrus.SetOutput(os.Stdout)
	api.Start()

}
