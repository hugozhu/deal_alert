package main

import (
	"flag"
	"fmt"
	"github.com/awsong/go-darts"
	"github.com/hugozhu/log4go"
	"io/ioutil"
	"os"
	"sqlite"
	"strings"
	"time"
	"unicode"
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
	log.DebugEnabled = *EnableDebug

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
		log.Info("[Enable debuging]")
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
				log.Debug(line)
				result := find_keywords(dict, line)
				if len(result) > 0 { //matched
					message := ""
					for k, _ := range result {
						var users []weibo.UserKeyword
						db.Query(&users, "select weibo_uid, keyword from user_keyword where keyword like ?", k)
						if users == nil || len(users) < 1 {
							continue
						}
						message = message + fmt.Sprintf("#%s# ", k)
						for _, u := range users {
							weibo_user := sina.UsersShow(u.WeiboUid)
							if weibo_user != nil {
								message = message + "@" + weibo_user.Screen_name + " "
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

func find_keywords(dict darts.Darts, line string) map[string]int {
	arr := []rune(strings.ToUpper(line))
	result := make(map[string]int)
	for i := 0; i < len(arr); i++ {
		offset := i
		c := arr[offset]
		if unicode.IsSpace(c) || unicode.IsPunct(c) {
			continue
		}
		for pos := 1; offset+pos < len(arr); pos++ {
			c := arr[offset+pos-1]
			if unicode.IsPunct(c) {
				break
			}
			// log.Info(string(arr[offset : offset+pos]))
			exist, results := dict.CommonPrefixSearch(arr[offset:offset+pos], 0)
			if len(results) > 0 {
				key := string(arr[offset : offset+pos])
				result[key] = result[key] + 1
				offset = offset + pos - 1
			} else if !exist {
				break
			}
		}
	}
	return result
}

func readToken() string {
	data, err := ioutil.ReadFile(os.Getenv("PWD") + "/token")
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return string(data[:32])
}
