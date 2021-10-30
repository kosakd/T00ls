package PutApi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

//判断是使用server酱还是推送加
func GetSckey(msg string, title string, sendkey string) string {
	var url_key string
	if strings.Contains(sendkey, "SCU") {
		url_key = "https://sc.ftqq.com/" + url.QueryEscape(sendkey) + ".send?text=" + url.PathEscape(title) + "&desp=" + url.QueryEscape(msg)
		return url_key
	} else if strings.Contains(sendkey, "SCT") {
		url_key = "https://sctapi.ftqq.com/" + url.QueryEscape(sendkey) + ".send?title=" + url.PathEscape(title) + "&desp=" + url.QueryEscape(msg)
		return url_key
	} else if len(sendkey) == 32 {
		url_key = "http://pushplus.hxtrip.com/send?token=" + url.QueryEscape(sendkey) + "&title=" + url.PathEscape(title) + "&content=" + url.QueryEscape(msg) + "&template=html"
		return url_key
	} else {
		panic("你的KEY设置有误请检查后，从新设置！！！")
	}
}

// 方糖推送
func Push(msg string, title string, sendkey string) {
	url_1 := GetSckey(msg, title, sendkey)
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
	} else if strings.Contains(string(body), "请求成功") {
		fmt.Println("Puts+发送成功")
	} else {
		fmt.Println("server酱发送失败")
	}
}
