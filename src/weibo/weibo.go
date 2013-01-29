package weibo

import (
	"encoding/json"
	"github.com/hugozhu/log4go"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const BaseURL = "https://api.weibo.com/2"

var log = log4go.New(os.Stdout)

func SetDebugEnabled(enable *bool) {
	log.DebugEnabled = enable
}

type Sina struct {
	AccessToken string
}

type UserKeyword struct {
	Keyword       string
	WeiboUid      int64
	WeiboUsername string
	Id            int64
	Frequence     int64
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
	Statuses []*WeiboPost
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

type WeiboComment struct {
	Id     int64
	Text   string
	Source string
	Mid    string
	User   *WeiboUser
	Status *WeiboPost
}

func (s *Sina) TimeLine(uid int64, since_id int64, count int) []*WeiboPost {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(uid, 10))
	params.Set("since_id", strconv.FormatInt(since_id, 10))
	params.Set("count", strconv.Itoa(count))
	var posts WeiboPosts
	if s.GET("/statuses/user_timeline.json", params, &posts) {
		return posts.Statuses
	}
	return nil
}

func (s *Sina) UsersShow(uid int64) *WeiboUser {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(uid, 10))
	var v WeiboUser
	if s.GET("/users/show.json", params, &v) {
		return &v
	}
	return nil
}

func (s *Sina) CommentsCreate(id int64, comment string) *WeiboComment {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(id, 10))
	params.Set("comment", comment)
	var v WeiboComment
	if s.POST("/comments/create.json", params, &v) {
		return &v
	}
	return nil
}

func (s *Sina) StatusesRepost(id int64, status string) *WeiboPost {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(id, 10))
	params.Set("status", status)
	params.Set("is_comment", "1")
	var v WeiboPost
	if s.POST("/statuses/repost.json", params, &v) {
		return &v
	}
	return nil
}

func (s *Sina) weiboApi(method string, base string, query url.Values, v interface{}) bool {
	url1 := BaseURL + base
	query.Set("access_token", s.AccessToken)

	var resp *http.Response
	var err error
	if method == "POST" {
		resp, err = http.PostForm(url1, query)
	} else {
		url1 += "?" + query.Encode()
		resp, err = http.Get(url1)
	}
	if err != nil {
		log.Error("fetch url %s %s", url1, err)
		panic(err)
	}
	defer resp.Body.Close()

	log.Debug(url1)

	if resp.StatusCode == 200 {
		d := json.NewDecoder(resp.Body)
		err = d.Decode(&v)
		if err != nil {
			panic(err)
		}
		return true
	} else {
		bytes, _ := ioutil.ReadAll(resp.Body)
		log.Error("Weibo API Error: " + string(bytes))
	}
	return false
}

func (s *Sina) GET(base string, query url.Values, v interface{}) bool {
	return s.weiboApi("GET", base, query, v)
}

func (s *Sina) POST(base string, query url.Values, v interface{}) bool {
	return s.weiboApi("POST", base, query, v)
}
