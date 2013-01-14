package main

import (
	"github.com/hugozhu/log4go"
	"os"
	"weibo"
)

var log = log4go.New(os.Stdout)

var sina *weibo.Sina

func init() {
	sina = &weibo.Sina{
		AccessToken: "2.008TkTLDIQdqsD4bbfd082cchG3E9E",
	}
}

func main() {
	log.Info("Hello")
	sina.TimeLine(2132734472, 0, 50)
}
