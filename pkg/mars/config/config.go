package config

import (
	"io"
	"os"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	MarsApiAddr   string `env:"MARS_API_ADDR" envDefault:":8080"`
	HashApiAddr   string `env:"HASH_API_ADDR" envDefault:"http://localhost:8082"`
	HashApiUser   string `env:"HASH_API_USER" envDefault:"mars"`
	HashApiPass   string `env:"HASH_API_PASS" envDefault:"test"`
	SmsCenterAddr string `env:"SMS_CENTER_ADDR" envDefault:"localhost:8082"`
	SmsCenterUser string `env:"SMS_CENTER_USER" envDefault:"mars"`
	SmsCenterPass string `env:"SMS_CENTER_PASS" envDefault:"test"`
	FtpAddr       string `env:"MARS_SFTP_ADDR" envDefault:"localhost:2222"`
	FtpUser       string `env:"MARS_SFTP_USER" envDefault:"mars"`
	FtpPass       string `env:"MARS_SFTP_PASS" envDefault:"test"`
	KafkaAddr     string `env:"MARS_QUEUE_ADDR" envDefault:"localhost:9093"`
	KafkaTopic    string `env:"MARS_QUEUE_TOPIC" envDefault:"mars_topic"`
	Log           *logrus.Logger
}

func NewConfig() *Config {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(io.MultiWriter(
		os.Stdout,
		&lumberjack.Logger{
			Filename:  "log/mars.log",
			MaxSize:   10,
			MaxAge:    3 * 365,
			LocalTime: true,
			Compress:  true,
		}))
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	c := &Config{}
	err := env.Parse(c)
	if err != nil {
		log.Fatal("config: unable to parse config from env: ", err)
	}

	c.Log = log

	return c
}
