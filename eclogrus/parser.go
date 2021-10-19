package eclogrus

import "github.com/sirupsen/logrus"

type Level logrus.Level

func (l *Level) Decode(value string) error {
	val, err := logrus.ParseLevel(value)
	*l = Level(val)
	return err
}
