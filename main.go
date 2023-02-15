package main

import (
	"Qscan/scan"
	"Qscan/spiderfinger"
	"os"
)

func main() {
	spiderfinger.Banner()
	scan.Scanurl(os.Args)
}
