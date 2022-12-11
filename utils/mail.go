package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"strconv"
	"strings"
)

type MailConfig struct {
	Host     string `key:"host" default:""`
	Port     string `key:"port" default:""`
	port     int
	User     string `key:"user" default:""`
	Password string `key:"password" default:""`
	From     string `key:"from" default:""`
	To       string `key:"to" default:""`
	toList   []string
}

var mailConfig MailConfig

func init() {
	ParseFromConfigFile("mail", &mailConfig)
	if mailConfig.Host != "" {
		var err error
		if mailConfig.port, err = strconv.Atoi(mailConfig.Port); err != nil {
			panic(fmt.Errorf("parse mail port err: %v", err))
		}
		mailConfig.toList = strings.Split(mailConfig.To, ",")
	}
}

func SendMail(subject string, body string) error {
	if mailConfig.Host == "" {
		logrus.Info("mail not configured, do nothing")
		return nil
	}
	dail := gomail.NewDialer(mailConfig.Host, mailConfig.port, mailConfig.User, mailConfig.Password)
	msg := gomail.NewMessage()
	msg.SetBody("text/html", body)
	msg.SetHeader("From", mailConfig.From)
	msg.SetHeader("To", mailConfig.toList...)
	msg.SetHeader("Subject", subject)
	return dail.DialAndSend(msg)
}
