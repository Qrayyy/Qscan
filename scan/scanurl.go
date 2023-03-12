package scan

import (
	"Qscan/http"
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
		//if Confluence.CVE_2022_26134(f.Url) {
		//	fmt.Println("[+]存在Confluence命令执行漏洞(CVE_2022_26134)")
		//} else {
		//	fmt.Println("[-]未发现漏洞")
		//}

		//var respBody string

		//respBody, err := Get(f.Url)
		//if err != nil {
		//	fmt.Println(err)
		//}

		links := spiderfinger.Spiderlinks(f.Url)
		fmt.Println(links)

		for _, link := range links {
			resp, err := http.Get(link, nil)
			if err != nil {
				fmt.Println(err)
			}
			var keyword string
			var flag int
			result, keyword, err := FingerRules.Matchkeyword("E:\\software\\GoLand 2021.2.2\\Projects\\Qscan\\scan\\FingerRules\\FingerRules.json", resp.Body)
			if result == true {
				//调用指定指纹poc
				flag = 1
			} else {
				//调用所有poc
			}
		}

		//fmt.Printf("[+]响应body:%s\n", respBody)
	}
}
