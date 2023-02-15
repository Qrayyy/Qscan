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

	//接受爬虫返回的指纹信息
	for finger := range spider.Result {
		fmt.Println(finger)
	}
}
