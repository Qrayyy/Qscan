package spiderfinger

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Finger struct {
	Url         string
	Title       string
	Server      string
	XPoweredBy  string
	ContentType string
}

type Spider struct {
	Result chan Finger
}

func (s *Spider) Runspider(args []string) {
	url := flag.String("u", "", "url")
	filename := flag.String("f", "", "filename")
	flag.Parse()
	if *url == "" && *filename == "" {
		fmt.Println("请输入url或file")
		return
	}
	var f Finger
	if *url != "" {
		f = s.SpiderUrl(*url)
		s.Result <- f
	}

	if *filename != "" {
		file, err := os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			f = s.SpiderUrl(scanner.Text())
			s.Result <- f
		}
		defer file.Close()
	}
	close(s.Result)
}

func (s *Spider) SpiderUrl(url string) Finger {
	// 创建一个colly实例
	c := colly.NewCollector()

	// 设置请求头
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	})

	// 获取网站指纹信息
	var finger Finger

	finger.Url = url
	c.OnResponse(func(r *colly.Response) {
		finger.Server = r.Headers.Get("Server")
		finger.XPoweredBy = r.Headers.Get("X-Powered-By")
		finger.ContentType = r.Headers.Get("Content-Type")
	})

	// 获取title和h1标签
	c.OnHTML("title", func(e *colly.HTMLElement) {
		finger.Title = e.Text
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(url)

	if err != nil {
		log.Fatalf("访问 %s 失败: %v", url, err)
	}

	return finger
}

func Spiderlinks(Url string) []string {
	// 创建一个colly实例
	c := colly.NewCollector()

	// 定义一个接口 URL 的数组
	var Links []string

	// 在请求发送之前，修改请求体
	c.OnRequest(func(r *colly.Request) {
		// 获取请求 URL 的主机名
		u, err := url.Parse(r.URL.String())
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}

		// 将请求体修改为 GET 请求，并带上 referer 和 origin 头部
		r.Method = http.MethodGet
		r.Headers.Set("Referer", r.URL.String())
		r.Headers.Set("Origin", u.Scheme+"://"+u.Host)
	})

	// 在响应接收之后，解析页面中的链接，并将它们存入数组中
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
			Links = append(Links, link)
		}
	})

	// 在访问页面时发生错误时，返回错误信息
	c.OnError(func(r *colly.Response, err error) {
		fmt.Errorf("failed to visit %s: %v", r.Request.URL, err)
	})

	// 访问给定的 URL
	if err := c.Visit(Url); err != nil {
		return nil
	}

	return Links
}

func (s *Spider) SpiderUrls(urls []string) {
	for _, u := range urls {
		go s.SpiderUrl(u)
	}
}
