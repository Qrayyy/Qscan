package scan

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

func spider(url string) string{
	// 创建一个colly实例
	c := colly.NewCollector()

	// 设置请求头
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	})

	// 获取网站指纹信息
	var Server string
	var XPoweredBy string
	var ContentType string
	c.OnResponse(func(r *colly.Response) {
		Server = r.Headers.Get("Server")
		XPoweredBy = r.Headers.Get("X-Powered-By")
		ContentType = r.Headers.Get("Content-Type")
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// 获取title和h1标签
	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("Title:", e.Text)
	})
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		fmt.Println("h1:", e.Text)
	})

	result := ""
	err := c.Visit(url)

	if err != nil{
		result = err.Error()
	}

	return result
}

