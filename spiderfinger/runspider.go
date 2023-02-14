package spiderfinger

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func Runspider(args []string) Finger {
	url := flag.String("u", "", "url")
	filename := flag.String("f", "", "filename")
	flag.Parse()

	if *url == "" && *filename == "" {
		fmt.Println("请输入url或file")
		return Finger{}
	}

	var finger Finger
	if *url != "" {
		finger = SpiderUrl(*url)
	}

	if *filename != "" {
		file, err := os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			finger = spider(scanner.Text())
		}
	}
	return finger
}
