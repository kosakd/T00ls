package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
	"tools/GetAllArticles"
	"tools/PutApi"
)

const (
	action      = "login"
	username    = "kosakd"                                 //用户名
	password    = "e10adc3949ba59abbe56e057f20f883e"       //密码md5 32位
	questionid  = "1"                                      //安全问题ID，默认0为未设置
	answer      = "kosakd"                                 //安全问题答案
	Sendkey     = "SCUxxxxxxxxxxxxxxxxxxx"                 //Server酱sendkey
	Url         = "https://www.t00ls.cc/All-articles.json" //api的url
	ToolsUrl    = "https://www.t00ls.cc"                   //tools的url
	Pages       = 10                                       //获取的页数
	PageNumbers = 10                                       //每页获取多少条数据
)

type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	Formhash   string `json:"formhash"`
	Mark       string `json:"mark"`
	Cookie     string
	Signsubmit string
}

var r Response

func main() {
	a := new(GetAllArticles.Articles)
	a.Url = Url
	a.ToolsUrl = ToolsUrl
	a.Pages = Pages
	a.PageNumbers = PageNumbers
	a.Sendkey = Sendkey
	a.Add_articles = make(chan int, a.Pages*a.PageNumbers)
	//此线程是死循环不用管它，它负责从Add_articles管道中取出文章位置，然后调用推送api，推送到手机中
	go a.Put_articles()
	//此线程是死循环不用管它，它负责从Add_articles管道中写入新增加文章位置
	GetAllArticles.Articles_list_All.Articles = make([]GetAllArticles.Articles_add, a.Pages*a.PageNumbers)
	go func() {
		for {
			a.Get_All_articles()
			time.Sleep(time.Minute * 15)
		}
	}()

	for {
		getCookie, _ := cookiejar.New(nil)
		client := &http.Client{Jar: getCookie}
		resp, _ := client.PostForm("https://www.t00ls.cc/login.json", url.Values{"action": {action}, "username": {username}, "password": {password}, "questionid": {questionid}, "answer": {answer}})
		json.NewDecoder(resp.Body).Decode(&r)
		resp.Body.Close()
		if r.Status != "success" {
			fmt.Println("登陆失败，一小时后重试。")
			time.Sleep(time.Hour)
			main()
		}
		r.Signsubmit = "true"
		ajaxsign(r, client)
		domainsearch(r, client, getdomain())
		time.Sleep(time.Hour)
		ajaxsign(r, client)
		time.Sleep(time.Hour * 24)

	}

}

// t00ls签到
func ajaxsign(r Response, client *http.Client) {
	resp, err := client.PostForm("https://www.t00ls.cc/ajax-sign.json", url.Values{"signsubmit": {r.Signsubmit}, "formhash": {r.Formhash}})
	if err != nil {
		fmt.Println("url访问失败：", err)
	}
	defer resp.Body.Close()
	var sign Response
	json.NewDecoder(resp.Body).Decode(&sign)
	if sign.Status == "success" {
		fmt.Println("1.签到成功")
		PutApi.Push(time.Now().Format("2006/01/02 15:04")+" 1.签到成功", "tools签到", Sendkey)
	} else if sign.Message == "alreadysign" {
		fmt.Println("2.今日已完成签到。")
		PutApi.Push(time.Now().Format("2006/01/02 15:04")+" 2.今日已完成签到", "tools签到", Sendkey)
	} else {
		fmt.Println("签到失败，1小时后重试。")
		PutApi.Push(time.Now().Format("2006/01/02 15:04")+" 3.签到失败，", "tools签到，1小时后重试或手动签到", Sendkey)
		time.Sleep(time.Hour)
		ajaxsign(r, client)
	}
}

// 获取最新备案域名，去重。
func getdomain() []string {
	temp := make([]string, 0)
	res := []string{}
	for i := 1; i <= 10; i++ {
		resp, err := http.Get(fmt.Sprintf("http://www.beianw.com/home/index/%d", i))
		if err != nil {
			fmt.Println("url获取失败", err)
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("url访问失败", err)
			continue
		}
		varhrefRegexp := regexp.MustCompile("\\w{0,62}\\.com")
		match := varhrefRegexp.FindAllString(string(body), -1)
		temp = append(temp, match...)
	}
	for i := range temp {
		flag := true
		for j := range res {
			if temp[i] == res[j] {
				flag = false
				break
			}
		}
		if flag {
			res = append(res, temp[i])
		}
	}
	return res
}

// 查询域名并查询tubi获取日志，如果包含域名则查询成功。
func domainsearch(r Response, client *http.Client, res []string) {
OneDomainsearch:
	for i := 2; i < len(res); i++ {
		client.PostForm("https://www.t00ls.cc/domain.html", url.Values{"domain": {res[i]}, "formhash": {r.Formhash}, "querydomainsubmit": {"%E6%9F%A5%E8%AF%A2"}})
		tubilog, err := client.Get("https://www.t00ls.cc/members-tubilog.json")
		if err != nil {
			fmt.Println("url查询失败：", err)
			continue OneDomainsearch
		}
		defer tubilog.Body.Close()
		body, err := io.ReadAll(tubilog.Body)
		if err != nil {
			fmt.Println("body获取失败：", err)
			continue OneDomainsearch
		}
		if strings.Contains(string(body), res[i]) {
			fmt.Printf("%s 域名查询成功，Tubi Get！", res[i])
			PutApi.Push(time.Unix(time.Now().Unix(), 0).UTC().Add(8*time.Hour).Format("2006-01-02 15:04:05")+res[i]+"域名查询成功", "tools域名查询", Sendkey)
			break OneDomainsearch
		}
	}
}
