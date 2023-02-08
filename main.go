package main

import (
	"fmt"

	"Qscan/scan"
	"flag"
	"os"
)

func main()  {
	url := flag.String("u","","单个url")
	flie := flag.String("f","","文件")
	flag.Parse()

	if *url != "" {
		args := []string{"-u", *url}
		if err := scan.spider(args); err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
	}else if *file != ""{
		args := []string{"-u", *url}
		if err := scan.spider(args); err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
	}else {
		fmt.Println("请输入url或文件")
		os.Exit(1)
	}
}