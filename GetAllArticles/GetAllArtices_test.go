package GetAllArticles

import (
	"fmt"
	"testing"
	"time"
)

func TestGetAllArticles(t *testing.T) {
	fmt.Println("test")
	// Get_all_articles("https://www.t00ls.cc/All-articles.json", "2")
	fmt.Println("==========")
	a_s := make([]string, 100)
	fmt.Println(a_s[1])
	fmt.Println("==========")

	a := new(Articles)
	a.Url = "https://www.t00ls.cc/All-articles.json"
	a.ToolsUrl = "https://www.t00ls.cc"
	a.Pages = 10
	a.PageNumbers = 10
	a.Add_articles = make(chan int, a.Pages*a.PageNumbers)
	//次线程是死循环不用管它，它负责从Add_articles管道中取出文章位置，然后调用推送api，推送到手机中
	go a.Put_articles()
	Articles_list_All.Articles = make([]Articles_add, a.Pages*a.PageNumbers)
	// a.Get_All_articles()
	// ch := make(chan int, a.Pages)
	// a.Get_one_articles(ch, a.Pages)
	// a.Get_All_articles()

	// ch := make(chan int, a.Pages)
	// for i := 0; i < a.Pages; i++ {
	// 	page_i := i + 1
	// 	a.Get_one_articles(ch, page_i)
	// }

	// ch := make(chan int, a.Pages)
	// a.Get_one_articles(ch, 4)
	// a.Get_one_articles(ch, a.Pages)
	for {
		a.Get_All_articles()
		time.Sleep(time.Minute * 15)
	}

	// J()
	// fmt.Printf("\n当前有多少篇文章：一共%d篇\n", Articles_list_All.AllNumbers)
	// fmt.Printf("当前文章最新一篇位置再：第%d位\n", Articles_list_All.Numbers)
	// fmt.Println("==========")
	// fmt.Println(len(Articles_list_All.Articles))
	// fmt.Println("==========")
	// s1 := strconv.Itoa(11)
	// fmt.Println(reflect.TypeOf(s1))
}
