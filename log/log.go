package log

import (
	"fmt"
	logrus "github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	fmt.Println("log init")
}

