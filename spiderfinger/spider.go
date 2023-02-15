package spiderfinger

import (
	"bufio"
	"flag"
	"fmt"
	wappalyzer "github.com/projectdiscovery/wappalyzergo"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Finger struct {
	Url  string
	Apps []string
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
	// 获取网站指纹信息
	var finger Finger

	finger.Url = url

	resp, err := http.DefaultClient.Get(url)
	data, err := ioutil.ReadAll(resp.Body)
	client, err := wappalyzer.New()
	if err != nil {
		log.Printf("wappalyzergo.New() failed: %s", err)
	} else {
		apps, err := client.Fingerprint(resp.Header, data)
		if err != nil {
			log.Printf("client.Fingerprint(%s) failed: %s", url, err)
		} else {
			finger.Apps = apps
		}
	}

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
