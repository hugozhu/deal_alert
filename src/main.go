package main

import (
	"github.com/hugozhu/log4go"
	"os"
	"sqlite"
	"time"
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
	log.Info(time.Now())
	var weibo_list []weibo.Weibo
	sqlite.Run(db_file, func(db *sqlite.DB) {
		db.Query(&weibo_list, "select * from weibo")
		post_chan := make(chan bool, len(weibo_list))
		for _, w := range weibo_list {
			go func(w weibo.Weibo) {
				last_id := w.LastId
				posts := sina.TimeLine(w.WeiboId, w.LastId, 50)
				for _, post := range posts {
					a, b := db.Execute("insert into queue (post_id, url,text, weibo_id,created) values (?,?,?,?,?)",
						post.Id, "", post.Text, post.User.Id, time.Now().Unix())
					log.Info(a, b)
					if post.Id > last_id {
						last_id = post.Id
					}
				}
				db.Execute("update weibo set last_id=? where id=?", last_id, w.Id)
				post_chan <- true
			}(w)
		}
		for i := 0; i < len(weibo_list); i++ {
			<-post_chan
		}
	})
}
