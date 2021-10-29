package PutApi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// 方糖推送
func Push(msg string, title string, sendkey string) {
	url_1 := "https://sc.ftqq.com/" + sendkey + ".send?text=" + url.QueryEscape(title) + "&desp=" + url.QueryEscape(msg)
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
