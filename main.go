package main

import (
	"Qscan/spiderfinger"
	"os"
)

func main() {
	spiderfinger.Banner()
	spider := &spiderfinger.Spider{
		Result: make(chan spiderfinger.Finger),
	}
	spider.Runspider(os.Args)
}
