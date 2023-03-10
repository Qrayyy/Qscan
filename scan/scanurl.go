package scan

import (
	"Qscan/scan/FingerRules"
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
		fmt.Println("\n[+]Start scan:")
		fmt.Printf("[+]%s[%s]\n", f.Url, f.Title)
		//var respBody string

		//respBody, err := Get(f.Url)
		//if err != nil {
		//	fmt.Println(err)
		//}

		links := spiderfinger.Spiderlinks(f.Url)
		fmt.Println(links)
		for _, link := range links {
			linkBody, err := Get(link)
			if err != nil {
				fmt.Println(err)
			}
			result, err := FingerRules.Matchkeyword("E:\\software\\GoLand 2021.2.2\\Projects\\Qscan\\scan\\FingerRules\\FingerRules.json", linkBody)
			if result == true {
				//调用poc
			} else {

			}
		}

		//fmt.Printf("[+]响应body:%s\n", respBody)
	}
}
