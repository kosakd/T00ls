package PutApi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

//判断SCKEY是新的还是旧的
func GetSckey(Sendkey string) string {
	if strings.Contains(Sendkey, "SCU") {
		return "text"
	} else if strings.Contains(Sendkey, "SCT") {
		return "title"
	} else {
		panic("server酱的KEY设置有误请检查后，从新设置！！！")
	}
}

// 方糖推送
func Push(msg string, title string, sendkey string, title_t string) {
	url_1 := "https://sc.ftqq.com/" + sendkey + ".send?" + url.QueryEscape(title_t) + "=" + url.QueryEscape(title) + "&desp=" + url.QueryEscape(msg)
	fmt.Println(url_1)
	resp, err := http.Get(url_1)
	if err != nil {
		fmt.Println("获取url失败: ", err)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(body))
	if strings.Contains(string(body), "success") {
		fmt.Println("server酱发送成功")
	} else {
		fmt.Println("server酱发送失败")
	}
}
