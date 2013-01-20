package main

import (
	"fmt"
	darts "github.com/awsong/go-darts"
	"github.com/hugozhu/log4go"
	"os"
	"sqlite"
	"strings"
	"time"
	"unicode"
	"weibo"
)

type WeiboPost struct {
	Id      int64
	Text    string
	WeiboId int64
	Created int
	PostId  int64
}

var DB_FILE = "data/deal_alert.db"

var log = log4go.New(os.Stdout)

func main() {
	sina := &weibo.Sina{
		AccessToken: "2.008TkTLDIQdqsD4bbfd082cchG3E9E",
	}

	dict, err := darts.Load("data/deals.lib")
	if err != nil {
		panic(err)
	}

	sqlite.Run(DB_FILE, func(db *sqlite.DB) {
		var posts []WeiboPost
		db.Query(&posts, "select * from queue where id > ?", 2636)
		for _, post := range posts {
			result := find_keywords(dict, line)
			if len(result) > 0 {
				// log.Info(post.Id, result, post.Text)
				for k, _ := range result {
					var users []weibo.UserKeyword
					db.Query(&users, "select weibo_uid, keyword from user_keyword where keyword like ?", k)
					message := fmt.Sprintf("#%s#", k)
					for _, u := range users {
						weibo_user := sina.UsersShow(u.WeiboUid)
						if weibo_user != nil {
							message = message + " @" + weibo_user.Screen_name
						}
					}
					// r := sina.StatusesRepost(post.PostId, message)
					r := sina.CommentsCreate(post.PostId, message)
					time.Sleep(15 * time.Second)
					if r != nil {
						log.Info("success alert:" + message)
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
		for pos := 2; offset+pos < len(arr); pos++ {
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
