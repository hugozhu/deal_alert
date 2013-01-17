package main

import (
	"fmt"
	mmsego "github.com/awsong/MMSEGO"
	darts "github.com/awsong/go-darts"
	"log"
	"strings"
)

func main() {
	_, err := darts.Import("data/deals.txt", "data/deals.lib", true)
	if err != nil {
		log.Fatal(err)
	}
	s := new(mmsego.Segmenter)
	s.Init("data/deals.lib")
	line := "膳魔师副牌：Elmundo Gap 艾蒙多 不锈钢保温壶 1L apple applink 99元（2件5折，折合49.5元/个，限华东、华北、东北，另有泰福高"
	line = " " + strings.ToUpper(line)
	offset := 0
	takeWord := func(off int, length int) {
		fmt.Printf("%s ", string(line[off-offset:off-offset+length]))
	}
	s.Mmseg(line[0:], offset, takeWord, nil, true)
	log.Println("\ndone")
}
