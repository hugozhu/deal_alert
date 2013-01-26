package main

import (
	"fmt"
	"log"
	"sqlite"
	"strings"
	"weibo"
)

var DB_FILE = "data/deal_alert.db"

func main() {
	log.Println("generating new dict...")
	sqlite.Run(DB_FILE, func(db *sqlite.DB) {
		var keywords []weibo.UserKeyword
		db.Query(&keywords, "select keyword, count(id) as frequence from user_keyword group by keyword")
		for _, keyword := range keywords {
			line := keyword.Keyword
			for _, v := range strings.Split(line, " ") {
				fmt.Printf("%s\t%d\n", strings.ToUpper(v), keyword.Frequence)
			}
		}
	})
}
