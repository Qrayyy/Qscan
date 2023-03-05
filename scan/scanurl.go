package scan

import (
	"Qscan/spiderfinger"
	"fmt"
)

func Scanurl(args []string) {
	spider := &spiderfinger.Spider{
		Result: make(chan spiderfinger.Finger, 10),
	}
	//执行爬虫
	spider.Runspider(args)

	for f := range spider.Result {
		fmt.Println(f)
		var respBody string

		respBody, err := Get(f.Url)
		if err != nil {
			fmt.Println(err)
		}

		links := spiderfinger.Spiderlinks(f.Url)
		for _, link := range links {
			linkBody, err := Get(link)

		}
		fmt.Printf("\n响应body:%s\n", respBody)
	}
}
