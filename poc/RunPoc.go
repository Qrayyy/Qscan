package poc

import (
	"Qscan/poc/Confluence"
	"Qscan/poc/zabbix"
)

type Poc struct {
	vulnid  string
	Execute func(url string) (bool, map[string]string)
}

var pocmap = map[string][]Poc{
	"zabbix": {
		{
			vulnid:  "CVE_2022_23131",
			Execute: zabbix.CVE_2022_23131,
		},
	},
	"confluence": {
		{
			vulnid:  "CVE_2022_26134",
			Execute: Confluence.CVE_2022_26134,
		},
	},
	// 添加更多关键词和poc函数
}

func RunPoc(Url string, keyword string) (bool, []map[string]string) {
	var success bool
	//用来存放成功利用的漏洞info
	var vulns []map[string]string

	//用两个for循环处理是否匹配到关键词的几种情况
	for k, v := range pocmap {
		if keyword != "" && k != keyword {
			//如果有关键词但是不在指纹里，则不打poc
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
	} else {
		success = true
	}

	return success, vulns
}
