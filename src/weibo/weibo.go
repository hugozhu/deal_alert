package weibo

import (
	"encoding/json"
	"github.com/hugozhu/log4go"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const BaseURL = "https://api.weibo.com/2"

var log = log4go.New(os.Stdout)

type Sina struct {
	AccessToken string
}

type Weibo struct {
	Id        int64
	WeiboId   int64
	Status    int
	LastId    int64
	WeiboName string
	Created   int64
	Modified  int64
}

type WeiboPosts struct {
	Statuses []WeiboPost
}

type WeiboPost struct {
	Created_At              string
	Id                      int64
	Mid                     string
	Text                    string
	Source                  string
	Trucated                bool
	In_Reply_To_Status_Id   string
	In_Reply_To_Screen_Name string
	Thumbnail_Pic           string
	Bmiddle_Pic             string
	Original_Pic            string
	User                    *WeiboUser
	Retweeted_Status        *WeiboPost
}

type WeiboUser struct {
	Id                int64
	Screen_name       string
	Name              string
	Location          string
	Description       string
	Url               string
	Profile_Image_Url string
	Verified_Reason   string
}

func (s *Sina) TimeLine(uid int64, since_id int64, count int) []WeiboPost {
	params := url.Values{}
	params.Set("access_token", s.AccessToken)
	params.Set("uid", strconv.FormatInt(uid, 10))
	params.Set("since_id", strconv.FormatInt(since_id, 10))
	params.Set("count", strconv.Itoa(count))

	url := BaseURL + "/statuses/user_timeline.json?" + params.Encode()
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var posts WeiboPosts
	log.Info(url)
	if resp.StatusCode == 200 {
		d := json.NewDecoder(resp.Body)
		err = d.Decode(&posts)
		if err != nil {
			panic(err)
		}
		return posts.Statuses
	} else {
		log.Error("failed to fetch:" + url)
	}
	return nil
}
