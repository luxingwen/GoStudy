package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"code.google.com/p/mahonia"
)

func main() {
	address := "http://222.179.234.148/default6.aspx"
	//	username := "201306571123"
	//	passwd := "10086zzz"

	req, err := http.Get(address)
	if err != nil {
		fmt.Println("err :", err)
	}
	res, _ := ioutil.ReadAll(req.Body)
	dec := mahonia.NewDecoder("gb2312")
	ret := dec.ConvertString(string(res))
	//fmt.Println("学生", )

	validat := regexp.MustCompile("<input[^>]*name=\"__VIEWSTATE\"[^>]*value=\"([^\"]*)\"[^>]*>")
	ret1 := validat.FindString(ret)
	validat = regexp.MustCompile("value=\"([^\"]*)\"")
	ret2 := validat.FindString(ret1)
	validat = regexp.MustCompile("\"([^\"]*)\"")
	v := validat.FindString(ret2)
	vs := v[1 : len(v)-1]
	client := &http.Client{}
	from := fmt.Sprintf("__VIEWSTATE=%s&txtYhm=201306571123&txtMm=10086zzz&rblJs=%s&btnDl=%s&tnamXw=yhdl&tbtnsXw=yhdl|xwxsdl", string(vs), dec.ConvertString("学生"), dec.ConvertString("登陆"))
	body := strings.NewReader(dec.ConvertString(from))
	resp, err := http.NewRequest("POST", address, body)
	if err != nil {
		fmt.Println("err: ", err)
	}
	resp.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	resp.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp.Header.Set("Connection", "keep-alive")
	resp.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	resp.Header.Set("Host", "222.179.234.148")
	resp.Header.Set("Origin", "http://222.179.234.148")
	resp.Header.Set("Referer", "http://222.179.234.148/default6.aspx")
	resp.Header.Set("Pragma", "no-cache")
	resp.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:44.0) Gecko/20100101 Firefox/44.0")
	resp.Header.Set("Accept-Encoding", "gzip, deflate")
	fmt.Println("resp ", resp)
	re, err := client.Do(resp)
	if err != nil {
		fmt.Println("err: ", err)
	}
	res, _ = ioutil.ReadAll(re.Body)
	dec = mahonia.NewDecoder("gb2312")
	ret = dec.ConvertString(string(res))
	fmt.Println(ret)

}
