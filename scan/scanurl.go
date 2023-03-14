package scan

import (
	"Qscan/http"
	"Qscan/poc"
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

		for _, Url := range links {
			resp, err := http.Get(Url, nil)
			if err != nil {
				fmt.Println(err)
			}
			keyword, err := FingerRules.Matchkeyword("E:\\software\\GoLand 2021.2.2\\Projects\\Qscan\\scan\\FingerRules\\FingerRules.json", resp.Body)
			success, vulnsInfos := poc.RunPoc(Url, keyword)
			if success {
				for _, vulnInfo := range vulnsInfos {
					for k, v := range vulnInfo {
						fmt.Printf("%s: %s\n", k, v)
					}
				}
			} else {
				fmt.Println("[-]未检测出漏洞")
			}
		}

		//fmt.Printf("[+]响应body:%s\n", respBody)
	}
}
