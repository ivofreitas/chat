package log

import (
	"context"
	"github.com/ivofreitas/chat/pkg/config"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type key string

func (k key) String() string {
	return "key: " + string(k)
}

var (
	log  *logrus.Logger
	once sync.Once
)

func Init() {
	once.Do(func() {
		logConfig := config.GetEnv().Log
		log = logrus.New()
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
		if level, err := logrus.ParseLevel(logConfig.Level); err != nil {
			logrus.SetLevel(logrus.WarnLevel)
		} else {
			log.SetLevel(level)
		}
	})
}

func InitParams(ctx context.Context) context.Context {

	httpLog := new(HTTP)
	httpLog.Request = new(Request)
	httpLog.Response = new(Response)

	ctx = context.WithValue(ctx, HTTPKey, httpLog)

	return ctx
}

func NewEntry() *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"mutex": &sync.Mutex{},
		"type":  "json",
	})
}
