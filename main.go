package main

import (
	"fmt"
	"net/http"
	"regexp"

	// "os"
	"strconv"
)

//查询网页内容
func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}

	defer resp.Body.Close()

	//读取网页body内容
	buf := make([]byte, 1024*4)
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("resp read body err = ", err)
			break
		}

		result += string(buf[0:n]) //累加读取的内容

	}

	return
}

func SpiderPage(i int) {
	//明确爬取的url
	url := "http://www.pengfu.com/xiaohua_" + strconv.Itoa(i) + ".html"
	fmt.Printf("正在爬取 %d个网页:%s\n", i, url)

	//开始爬取网页内容
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("Http get err = ", err)
		return
	}

	// fmt.Println("r = ", result)
	//取，一个段子url的连接
	//解析表达式
	re := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?)"`)
	if re == nil {
		fmt.Println("regexp compile failed")
		return
	}
	//取关键信息
	joyUrls := re.FindAllStringSubmatch(result, -1)
	// fmt.Println("joyUrls = ", joyUrls)

	//取网址
	//第一个返回下标，第二个返回内存
	for _, data := range joyUrls {
		fmt.Println("url = ", data[1])
	}

}

func DoWork(start, end int) {

	fmt.Printf("准备爬取第%d到%d页的网址\n", start, end)
	for i := start; i <= end; i++ {
		SpiderPage(i)
	}

}

func main() {
	var start, end int
	fmt.Printf("请输入起始页（>=1）: ")
	fmt.Scan(&start)
	fmt.Printf("请输入终止页（>=起始页）: ")
	fmt.Scan(&end)

	DoWork(start, end)
}
