package search

import (
	"github.com/awsong/go-darts"
	"github.com/hugozhu/log4go"
	"os"
	"strings"
	"unicode"
)

var log = log4go.New(os.Stdout)

func SetDebugEnabled(enable bool) {
	log.DebugEnabled = &enable
}

func FindKeywords(dict darts.Darts, line string) map[string]int {
	arr := []rune(strings.ToUpper(line))
	result := make(map[string]int)
	for i := 0; i < len(arr); i++ {
		offset := i
		c := arr[offset]
		if unicode.IsSpace(c) {
			log.Debug("ignore space：", c)
			continue
		}
		for pos := 1; offset+pos < len(arr)+1; pos++ {
			// c := arr[offset+pos-1]
			// if unicode.IsPunct(c) {
			// 	break
			// }
			log.Debug(string(arr[offset : offset+pos]))
			exist, results := dict.CommonPrefixSearch(arr[offset:offset+pos], 0)
			if len(results) > 0 {
				not_whole_english_word := false
				if IsEnglishWord(arr[offset : offset+pos]) {
					if offset >= 1 && !IsEnglishEdge(arr[offset-1]) {
						log.Debug("skiped", string(arr[offset:offset+pos]))
						not_whole_english_word = true
					}
					if offset+pos < len(arr)+1 && !IsEnglishEdge(arr[offset+pos]) {
						log.Debug("skiped", string(arr[offset:offset+pos]))
						not_whole_english_word = true
					}
				}
				if !not_whole_english_word {
					key := string(arr[offset : offset+pos])
					result[key] = result[key] + 1
				}
				offset = offset + pos - 1
			} else if !exist {
				break
			}
		}
	}
	return result
}

func IsEnglishEdge(r rune) bool {
	if 'A' <= r && r <= 'Z' {
		return false
	}
	return true
}

func IsEnglishWord(arr []rune) bool {
	for _, r := range arr {
		if int32(r) > unicode.MaxLatin1 {
			return false
		}
	}
	return true
}
