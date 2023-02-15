package spiderfinger

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

type Finger struct {
	Url   string
	Title string
	//H1          string
	Server      string
	XPoweredBy  string
	ContentType string
}

type Spider struct {
	Result chan Finger
}

func (s *Spider) Runspider(args []string) {
	//s.Result = make(chan Finger, 10)
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
		fmt.Println(f)
		select {
		case s.Result <- f:
			fmt.Println(s)
		default:
			fmt.Println("管道已关闭")
		}
	}

	if *filename != "" {
		file, err := os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			f = s.SpiderUrl(scanner.Text())
			s.Result <- f
		}
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
		//finger.H1 = e.Text
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

func (s *Spider) SpiderUrls(urls []string) {
	for _, url := range urls {
		go s.SpiderUrl(url)
	}
}
