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
	}
}
