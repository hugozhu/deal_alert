package main

import (
	"flag"
	"fmt"
	"github.com/awsong/go-darts"
	"github.com/hugozhu/log4go"
	"io/ioutil"
	"os"
	"search"
	"sqlite"
	"time"
	"weibo"
)

var EnableDebug = flag.Bool("debug", false, "enable debug")
var IsTestMode = flag.Bool("test", false, "enable test mode")

var log = log4go.New(os.Stdout)

var sina *weibo.Sina
var DB_FILE = os.Getenv("PWD") + "/data/deal_alert.db"
var DICT_FILE = os.Getenv("PWD") + "/data/deals.lib"

var dict darts.Darts

func init() {
	sina = &weibo.Sina{
		AccessToken: readToken(),
	}
	log.DebugEnabled = EnableDebug
	weibo.SetDebugEnabled(EnableDebug)

	var err error
	dict, err = darts.Load("data/deals.lib")
	if err != nil {
		panic(err)
	}

	flag.Parse()
	if *IsTestMode {
		log.Info("[Test mode]")
	}
	if *EnableDebug {
		log.Debug("[Enable debuging]")
	}
}

func main() {
	var weibo_list []weibo.Weibo
	sqlite.Run(DB_FILE, func(db *sqlite.DB) {
		db.Query(&weibo_list, "select * from weibo")
		post_chan := make(chan []*weibo.WeiboPost, len(weibo_list))
		for _, w := range weibo_list {
			go func(w weibo.Weibo) {
				if *IsTestMode {
					w.LastId = 0
				}
				last_id := w.LastId
				posts := sina.TimeLine(w.WeiboId, w.LastId, 10)
				posts2 := []*weibo.WeiboPost{}
				for _, post := range posts {
					if post.Id <= w.LastId {
						//ignore 过期置顶贴
						continue
					}
					if post.Retweeted_Status != nil {
						post.Text = post.Text + "//" + post.Retweeted_Status.Text
						//ignore 转帖，重复推荐意义不大
						continue
					}

					// log.Info(post.Text)
					if !*IsTestMode {
						_, err := db.Execute("insert into queue (post_id, url,text, weibo_id,created) values (?,?,?,?,?)",
							post.Id, "", post.Text, post.User.Id, time.Now().Unix())
						if err != nil {
							log.Error("Failed to save ", post.Id, w.LastId, post)
						} else {
							posts2 = append(posts2, post)
						}
						if post.Id > last_id {
							last_id = post.Id
						}
					} else {
						posts2 = append(posts2, post)
					}
				}
				if !*IsTestMode {
					db.Execute("update weibo set last_id=? where id=?", last_id, w.Id)
				}
				post_chan <- posts2
			}(w)
		}
		for i := 0; i < len(weibo_list); i++ {
			posts := <-post_chan
			for _, post := range posts {
				line := post.Text
				result := search.FindKeywords(dict, line)
				if len(result) > 0 { //matched
					urls := search.FindUrls(line)
					if len(urls) > 0 {
						for _, u := range weibo.ExpandUrls(urls) {
							line = line + " " + u
						}
					}
					message := ""
					at_users := make(map[string]bool)
					for k, _ := range result {
						var users []weibo.UserKeyword
						db.Query(&users, "select weibo_uid, keyword from user_keyword where keyword like ?", k)
						if users == nil || len(users) < 1 {
							continue
						}

						message = message + fmt.Sprintf("#%s# ", k)
						for _, u := range users {
							weibo_user := sina.UsersShow(u.WeiboUid)
							if weibo_user != nil && at_users[weibo_user.Screen_name] != true {
								message = message + "@" + weibo_user.Screen_name + " "
								at_users[weibo_user.Screen_name] = true
							}
						}
					}
					if len(message) > 0 {
						if *IsTestMode {
							log.Info("success alert:", message, line)
						} else {
							r := sina.StatusesRepost(post.Id, message)
							//r := sina.CommentsCreate(post.Id, message)
							if r != nil {
								log.Info("success alert:" + message)
							}
							time.Sleep(3 * time.Second)
						}
					}
				}
			}
		}
	})
}

func readToken() string {
	data, err := ioutil.ReadFile(os.Getenv("PWD") + "/token")
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return string(data[:32])
}
