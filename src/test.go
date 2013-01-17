package main

import (
	"fmt"
	mmsego "github.com/awsong/MMSEGO"
	"sqlite"
	"strings"
)

type WeiboPost struct {
	Id      int64
	Text    string
	WeiboId int64
	Created int
	PostId  int64
}

var DB_FILE = "data/deal_alert.db"

func main() {
	s := new(mmsego.Segmenter)
	s.Init("data/deals.lib")

	sqlite.Run(DB_FILE, func(db *sqlite.DB) {
		var posts []WeiboPost
		db.Query(&posts, "select * from queue order by id asc")
		for _, post := range posts {
			line := strings.ToUpper(post.Text)
			offset := 0
			takeWord := func(off int, length int) {
				if length > 3 {
					fmt.Printf("%d %s\n", post.Id, string(line[off-offset:off-offset+length]))
				}
			}
			s.Mmseg(line[0:], offset, takeWord, nil, true)
		}
	})
}
