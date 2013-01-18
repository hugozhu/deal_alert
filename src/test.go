package main

import (
	darts "github.com/awsong/go-darts"
	"github.com/hugozhu/log4go"
	"os"
	"sqlite"
	"strings"
	"unicode"
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
	dict, err := darts.Load("data/deals.lib")
	if err != nil {
		panic(err)
	}

	sqlite.Run(DB_FILE, func(db *sqlite.DB) {
		var posts []WeiboPost
		db.Query(&posts, "select * from queue order by id asc limit 100 offset 0")
		for _, post := range posts {
			line := strings.ToUpper(post.Text)
			// log.Info(post.Id, line)
			result := find_keywords(dict, line)
			if len(result) > 0 {
				log.Info(post.Id, result, post.Text)
			}
		}
	})
}

func find_keywords(dict darts.Darts, line string) []string {
	arr := []rune(strings.ToUpper(line))
	result := []string{}
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
				result = append(result, string(arr[offset:offset+pos]))
				offset = offset + pos - 1
			} else if !exist {
				break
			}
		}
	}
	return result
}
