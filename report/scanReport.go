package report

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ScanReport(Url string, vulnInfo map[string]string) {

	info, _ := os.ReadFile("./report/vulnInfo.js")
	stringInfo := string(info)
	aInfo := []string{
		"\"plugin\":\"",
		"\"payload\":\"",
		"\"snapshot\":[[\"",
		"\"\"]]",
		"\"create_time\":",
		"\"addr\":\"",
		"\"url\":\"",
	}
	//用来存放传进来的漏洞信息
	bInfo := []string{"1", "1", "1", "1", "1", "1", "1"}
	var flag = 0
	for k := range vulnInfo {
		bInfo[flag] = vulnInfo[k]
		flag++
		if flag == 4 {
			t := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
			bInfo[flag] = t
			break
		}
	}
	bInfo[5] = Url
	bInfo[6] = Url

	for i := 0; i < 7; i++ {
		var a = aInfo[i]
		var b = bInfo[i]
		if i != 3 {
			var c = a + b
			stringInfo = strings.ReplaceAll(stringInfo, a, c)
		} else {
			var c = "\"" + b + "\"]]"
			stringInfo = strings.ReplaceAll(stringInfo, a, c)
		}
	}

	f, err := os.OpenFile("./report/vulnInfo.html", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("读取文件出错")
		return
	}
	defer f.Close()

	_, err = f.WriteString(stringInfo)
	if err != nil {
		return
	}
	fmt.Println("漏洞报告已生成!")
}
