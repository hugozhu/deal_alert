package main

import (
	"github.com/hugozhu/log4go"
	"os"
	"sqlite"
	"weibo"
)

var log = log4go.New(os.Stdout)

var sina *weibo.Sina

func init() {
	sina = &weibo.Sina{
		AccessToken: "2.008TkTLDIQdqsD4bbfd082cchG3E9E",
	}
}

var db_file = os.Getenv("PWD") + "/data/deal_alert.db"

func main() {
	sqlite.Run(db_file, func(db *sqlite.DB) {
		var v weibo.Weibo
		db.QueryResult(&v, "select * from weibo")
		log.Info(v)
	})

	weibo_ids := []int64{}
	complete_chan := make(chan bool, len(weibo_ids))

	for _, id := range weibo_ids {
		go func() {
			posts := sina.TimeLine(id, 0, 50)
			log.Info(id, len(posts))
			complete_chan <- true
			for _, post := range posts {
				log.Info(post.Id)
			}
		}()
	}
	for i := 0; i < len(weibo_ids); i++ {
		<-complete_chan
	}
}
