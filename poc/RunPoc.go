package poc

import (
	"Qscan/poc/Confluence"
	"Qscan/poc/nacos"
	"Qscan/poc/zabbix"
	"fmt"
	"github.com/fatih/color"
)

type Poc struct {
	vulnid  string
	Execute func(url string) (bool, map[string]string)
}

var pocmap = map[string][]Poc{
	"Zabbix SIA": {
		{
			vulnid:  "CVE_2022_23131",
			Execute: zabbix.CVE_2022_23131,
		},
	},
	"Atlassian Confluence": {
		{
			vulnid:  "CVE_2022_26134",
			Execute: Confluence.CVE_2022_26134,
		},
	},
	"Nacos": {
		{
			vulnid:  "CVE_2021_29441",
			Execute: nacos.CVE_2021_29441,
		},
	},
	// 添加更多关键词和poc函数
}

func RunPoc(Url string, keyword string) (bool, []map[string]string) {
	green := color.New(color.FgGreen).SprintFunc()

	var success bool
	//用来存放成功利用的漏洞info
	var vulns []map[string]string

	//用两个for循环处理匹配关键词的几种情况
	for k, v := range pocmap {
		if keyword != "" && k != keyword {
			//如果有关键词但是不在指纹里或无关键词，则不打poc
			continue
		}
		for _, poc := range v {
			success, info := poc.Execute(Url)
			if success {
				vulns = append(vulns, info)
			}
		}
	}

	//检验是否产生了漏洞，方便其他函数调用
	if len(vulns) == 0 {
		success = false
		fmt.Printf("%s\n\n", green("[-]未检测出漏洞"))
	} else {
		success = true
	}

	return success, vulns
}
