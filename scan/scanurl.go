package scan

import (
	"Qscan/http"
	"Qscan/poc"
	"Qscan/report"
	"Qscan/scan/FingerRules"
	"Qscan/spiderfinger"
	"fmt"
	"github.com/fatih/color"
)

func Scanurl(args []string) {
	spider := &spiderfinger.Spider{
		Result: make(chan spiderfinger.Finger, 10),
	}
	//打印日志颜色
	red := color.New(color.FgRed).SprintFunc()
	//执行爬虫
	spider.Runspider(args)
	fmt.Println("\n[+]Start scan:")
	for f := range spider.Result {
		fmt.Printf("[+]%s[%s]\n", f.Url, f.Title)

		firstUrl := f.Url
		resp, err := http.Get(firstUrl, nil)
		if err != nil {
			fmt.Println(err)
		}
		keyword, err := FingerRules.Matchkeyword("E:\\software\\GoLand 2021.2.2\\Projects\\Qscan\\scan\\FingerRules\\FingerRules.json", resp.Body)
		success, vulnsInfos := poc.RunPoc(firstUrl, keyword)
		if success {
			for _, vulnInfo := range vulnsInfos {
				fmt.Printf("%s\n", red("[+]漏洞信息如下："))
				for k, v := range vulnInfo {
					//打印日志
					fmt.Printf("%s: %s\n", k, v)
				}
				//导出报告
				report.ScanReport(firstUrl, vulnInfo)
				fmt.Println()
			}
		} else {
			links := spiderfinger.Spiderlinks(f.Url)
			//fmt.Println(links)

			for _, Url := range links {
				resp, err := http.Get(Url, nil)
				if err != nil {
					fmt.Println(err)
				}
				keyword, err = FingerRules.Matchkeyword("E:\\software\\GoLand 2021.2.2\\Projects\\Qscan\\scan\\FingerRules\\FingerRules.json", resp.Body)
				success, vulnsInfos = poc.RunPoc(Url, keyword)
				if success {
					for _, vulnInfo := range vulnsInfos {
						fmt.Println("[+]漏洞信息如下：")
						for k, v := range vulnInfo {
							//打印日志
							fmt.Printf("%s: %s\n", k, v)
						}
						//导出报告
						report.ScanReport(Url, vulnInfo)
						fmt.Println()
					}
				}
			}
		}
	}
}
