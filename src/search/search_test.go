package search

import (
	"github.com/awsong/go-darts"
	"os"
	"testing"
)

func TestEnglishWord(t *testing.T) {
	if !IsEnglishWord([]rune("hello")) {
		t.Error("hello is english word")
	}
	if IsEnglishWord([]rune("硬盘")) {
		t.Error("硬盘 is not english word")
	}
}

func TestSearch(t *testing.T) {
	SetDebugEnabled(false)
	dict, err := darts.Load(os.Getenv("PWD") + "/data/deals.lib")
	if err != nil {
		t.Fatal(err)
	}
	results := FindKeywords(dict, "CANON的相机 is a lvbag BURT'S BEES 儿童.GNC.硬盘")
	if results["CANON"] != 1 || results["BURT'S BEES"] != 1 || results["CANON"] != 1 ||
		results["硬盘"] != 1 || results["LV"] == 1 {
		t.Error("failed to find keyword")
	}
}

// func TestFindUrls(t *testing.T) {
// 	str := "$12.99 Newegg.com - 森海塞尔(Sennheiser)CX200 3.5mm Connector Canal Stereo Headphone 耳机http://t.cn/zYL7y4P（来自@北美省钱快报 iPhone客户端，下载链接: http://t.cn/zOkIU4D）"
// 	t.Log(ExpandUrlInfo(FindUrls(str), t))
// }
