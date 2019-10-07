package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

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

func SpiderOneJoy(url string) (title, content string, err error) {
	//开始爬取网页内容
	result, err1 := HttpGet(url)
	if err != nil {
		err = err1
		return
	}

	//取关键信息
	//取标题
	re1 := regexp.MustCompile(`<h1>(?s:(.*?))</h1>`)
	if re1 == nil {
		fmt.Println("regexp compile failed")
		return
	}
	//取内容
	tmpTitle := re1.FindAllStringSubmatch(result, 1) //最后参数为1，只过滤第一个
	for _, data := range tmpTitle {
		title = data[1]
		title = strings.Replace(title, "\t", "", -1) //将换行符替换为空字符
		break
	}

	re2 := regexp.MustCompile(`<div class="content-txt pt10">(?s:(.*?))<a id="prev" href="`)
	if re2 == nil {
		fmt.Println("regexp compile failed")
		return
	}
	//取内容
	tmpContent := re2.FindAllStringSubmatch(result, -1)
	for _, data := range tmpContent {
		content = data[1]
		content = strings.Replace(content, "\t", "", -1)
		content = strings.Replace(content, "\n", "", -1)
		content = strings.Replace(content, "\r", "", -1)
		break
	}

	return

}

func StoreJoyToFile(i int, fileTitle, fileContent []string) {
	//新建文件
	f, err := os.Create(strconv.Itoa(i) + ".txt")
	if err != nil {
		fmt.Println("os create err = ", err)
		return
	}
	defer f.Close()

	//写内存
	n := len(fileTitle)
	for i := 0; i < n; i++ {
		//写标题
		f.WriteString(fileTitle[i] + "\n")
		//写内容
		f.WriteString(fileContent[i] + "\n")

		f.WriteString("\n=======================================\n")
	}
}

func SpiderPage(i int, page chan int) {
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
	//取，<h1 class="dp-b"><a href= 一个段子url的连接
	//解析表达式
	re := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))"`)
	fmt.Println("re = ", re)
	if re == nil {
		fmt.Println("regexp compile failed")
		return
	}
	//取关键信息
	joyUrls := re.FindAllStringSubmatch(result, -1)
	// fmt.Println("joyUrls = ", joyUrls)

	fileTitle := make([]string, 0)
	fileContent := make([]string, 0)

	//取网址
	//第一个返回下标，第二个返回内存
	for _, data := range joyUrls {
		// fmt.Println("url = ", data[1])
		//爬取每一个笑话，每一个段子
		title, content, err := SpiderOneJoy(data[1])
		if err == nil {
			fmt.Println("Spider one joy err = ", err)
			continue
		}

		fileTitle = append(fileTitle, title)
		fileContent = append(fileContent, content)
	}

	// fmt.Println("fileTitle = ", fileTitle)
	// fmt.Println("fileContent = ", fileContent)

	//把内容写入到文件
	StoreJoyToFile(i, fileTitle, fileContent)

	page <- i

}

func DoWork(start, end int) {

	fmt.Printf("准备爬取第%d到%d页的网址\n", start, end)

	page := make(chan int)

	for i := start; i <= end; i++ {
		//爬取主页面
		go SpiderPage(i, page)
	}

	for i := start; i <= end; i++ {
		fmt.Printf("第%d个页面爬取完成\n", <-page)
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
