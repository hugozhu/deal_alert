package main

import (
	"fmt"
	"log"
	"sqlite"
	"strings"
)

type UserKeyword struct {
	Keyword   string
	Frequence int64
}

var DB_FILE = "data/deal_alert.db"

func main() {
	log.Println("generating new dict...")
	sqlite.Run(DB_FILE, func(db *sqlite.DB) {
		var keywords []UserKeyword
		db.Query(&keywords, "select keyword, count(id) as frequence from user_keyword group by keyword")
		for _, keyword := range keywords {
			fmt.Printf("%s\t%d\n", strings.ToUpper(keyword.Keyword), keyword.Frequence)
		}
	})
}
