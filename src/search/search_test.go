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
	results := FindKeywords(dict, "CANON的相机 is a lvbag BURT'S BEES 儿童GNC硬盘")
	if results["CANON"] != 1 || results["BURT'S BEES"] != 1 || results["CANON"] != 1 ||
		results["硬盘"] != 1 || results["LV"] == 1 {
		t.Error("failed to find keyword")
	}
}
