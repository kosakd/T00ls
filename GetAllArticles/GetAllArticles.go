package GetAllArticles

//当前包是获取文章的函数封装

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tools/PutApi"
)

//all_articles结构体是geturl是
type Articles struct {
	Url          string   //api的url
	Pages        int      //获取的页数
	PageNumbers  int      //每页获取多少条数据
	Add_articles chan int //记录新添加文章的位置的chan
	ToolsUrl     string   //tools的url
	Sendkey      string   //server酱的密钥
	Order        bool     //状态判断，判断是同过头部递增添加还是，通过尾部递减添加
}

//文章api的json结构体
type Articles_Response struct {
	Status       string `json:"status"`
	Articleslist []articleslist
}

type Articles_add struct {
	Subject string
	Message string
	Author  string
	Links   string
}

//所有记录所有文章的结构体
type Articles_list struct {
	Articles   []Articles_add //文章数据结构体
	Numbers    int            //最新一篇文文章的位置
	AllNumbers int            //当前总共有多少篇文章
}

type articleslist struct {
	Tid        string `json:"tid"`
	Subject    string `json:"subject"`
	Message    string `json:"message"`
	Dateline   string `json:"dateline"`
	Attachment string `json:"attachment"`
	Fid        string `json:"fid"`
	Author     string `json:"author"`
	Authorid   string `json:"authorid"`
	Views      string `json:"views"`
	Replies    string `json:"replies"`
	Pic        string `json:"pic"`
	Links      string `json:"links"`
	Avatar     string `json:"avatar"`
}

var Articles_list_All Articles_list

//获取当前一页文章的10篇文章title，返回一个新增加文章的links的list，
func (A *Articles) Get_one_articles(ch chan int, page int) []string {
	//先创建一个，文章结构体，用于添加数据
	var articles_one Articles_add

	all_articles := make([]string, 0)
	//声明client 参数为默认
	client := &http.Client{}
	//可以添加多个cookie
	cookie1 := &http.Cookie{Name: "page", Value: strconv.Itoa(page)}

	//声明要访问的url

	//钩爪请求请求
	reqest, err := http.NewRequest("GET", A.Url, nil)
	if err != nil {
		panic(err)
	}
	//添加header头中的ua头
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	//把创建cookie结构体放入请求体中
	reqest.AddCookie(cookie1)

	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Println("获取文章失败")
		panic(err)
	}
	defer response.Body.Close()

	//打印请求返回值
	body, _ := io.ReadAll(response.Body)
	//判断是否包含888这个int，此int是匿名的id，他妈的，好好的js返回string突然给爷整个int，干你大爷
	body_s := string(body)
	if strings.Contains(body_s, "888,") {
		old := "888,"
		new := `"888",`
		body_s = strings.Replace(body_s, old, new, -1)
	}

	//创建Articles_Response结构体，用于json的解析
	var Articles_json Articles_Response
	err = json.Unmarshal([]byte(body_s), &Articles_json)
	if err != nil {
		fmt.Println("json解析错误：", err)
		ch <- 0
		return nil
	}

	if Articles_json.Status != "success" {
		fmt.Println("文章接口出错：", err)
		ch <- 0
		return nil
	}
OuterLoopOne:
	for i := 0; i < A.PageNumbers; i++ {
		fmt.Println("haha")
		articles_one.Subject = Articles_json.Articleslist[i].Subject
		articles_one.Author = Articles_json.Articleslist[i].Author
		articles_one.Message = Articles_json.Articleslist[i].Message

		//用一个循环来判断当前的的文章是否已经收录了
		//这一块可以封装成一个函数，然后用go()去执行，把新增文章列表改为channel类型，这样就不怕添加文章重复了
		if Articles_list_All.AllNumbers == 0 {
			Articles_list_All.Articles[Articles_list_All.Numbers] = articles_one
			Articles_list_All.AllNumbers += 1
			continue OuterLoopOne
		}
		for i_1 := 0; i_1 < Articles_list_All.AllNumbers; i_1++ {
			if strings.Compare(articles_one.Links, Articles_list_All.Articles[i_1].Links) == 0 {
				fmt.Println("==========")
				fmt.Println(articles_one.Subject)
				fmt.Println("此片文章是重复的")
				continue OuterLoopOne
			}
		}
		fmt.Println("==========")
		fmt.Println(articles_one.Subject)
		fmt.Println("此片文章是新增加的")
		if Articles_list_All.AllNumbers == A.Pages*A.PageNumbers {
			if A.Order {
				Articles_list_All.Numbers += 1
				Articles_list_All.Articles[Articles_list_All.Numbers] = articles_one
				if Articles_list_All.Numbers == A.PageNumbers*A.Pages-1 {
					A.Order = false
				}
			} else {
				Articles_list_All.Numbers -= 1
				Articles_list_All.Articles[Articles_list_All.Numbers] = articles_one
				if Articles_list_All.Numbers == 0 {
					A.Order = true
				}
			}
		} else {
			Articles_list_All.Numbers += 1
			Articles_list_All.Articles[Articles_list_All.Numbers] = articles_one
			Articles_list_All.AllNumbers += 1
		}

		A.Add_articles <- Articles_list_All.Numbers
		//循环添加完毕，打印文章的标题和内容
		fmt.Println("==========")
		fmt.Printf("标题：%s,\n", Articles_list_All.Articles[Articles_list_All.Numbers].Subject)
		fmt.Printf("作者：%s,\n", Articles_list_All.Articles[Articles_list_All.Numbers].Author)
		fmt.Printf("内容：%s,\n", Articles_list_All.Articles[Articles_list_All.Numbers].Message)
		fmt.Printf("网址：%s,\n", Articles_list_All.Articles[Articles_list_All.Numbers].Links)
		fmt.Println("==========")

		//将获取到的文章放入一个
	}

	/*
		//通过创建接口来解析json
		var i interface{}
		json.Unmarshal(body, i)

		fmt.Println("打印接口解析的json")

		fmt.Println(i)
	*/
	//返回的状态码

	status := response.Status
	fmt.Println()
	fmt.Println("1111111111111111111111", status)
	fmt.Printf("当前是第%d次获取api\n", page)
	ch <- 1
	return all_articles
}

func (A *Articles) Get_All_articles() {
	ch := make(chan int, A.Pages)
	for i := 0; i < A.Pages; i++ {
		page_i := i + 1
		A.Get_one_articles(ch, page_i)
	}

}

func (A *Articles) Put_articles() {
	var add_numbers int = 0
Again:
	for {
		//如果想要前（A.PageNumbers*A.Pages）篇文章推送的话，，记得把这个if代码快和Again:给注释掉
		if add_numbers < A.PageNumbers*A.Pages {
			add_articles := <-A.Add_articles
			fmt.Println(add_articles)
			add_numbers += 1
			continue Again
		}
		add_articles := <-A.Add_articles
		fmt.Println(add_articles)
		test := []string{
			"标题：" + Articles_list_All.Articles[add_articles].Subject + ",\n",
			"作者：" + Articles_list_All.Articles[add_articles].Author + ",\n",
			"  " + Articles_list_All.Articles[add_articles].Message + ",\n",
			"网址：" + A.ToolsUrl + "/" + Articles_list_All.Articles[add_articles].Links + ",\n",
		}
		test_1 := strings.Join(test, "\n")
		//将新获取的文章信息通过推送api，发出去
		PutApi.Push(test_1, Articles_list_All.Articles[add_articles].Subject, A.Sendkey)
		//睡眠一秒防止，快速发包导致down包
		time.Sleep(time.Second * 5)

	}
}
