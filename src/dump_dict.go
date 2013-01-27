package main

import (
	"fmt"
	"log"
	"parse"
	"sqlite"
	"strings"
	"weibo"
)

var DB_FILE = "data/deal_alert.db"

func main() {
	log.Println("generating new dict...")
	sqlite.Run(DB_FILE, func(db *sqlite.DB) {
		var keywords []weibo.UserKeyword
		db.Query(&keywords, "select keyword, id from user_keyword")
		for _, keyword := range keywords {
			line := keyword.Keyword
			query, _ := parse.Parse(line)
			for _, k := range query.AllTerms() {
				fmt.Printf("%s\t%d\n", strings.ToUpper(k), keyword.Id)
			}
		}
	})
}
