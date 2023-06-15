package log

import "github.com/sirupsen/logrus"

type Logger struct {
	*logrus.Logger
	Pid int
}
