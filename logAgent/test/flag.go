package main

import (
	"flag"
	"fmt"
)

func main() {
	confile := flag.String("conf", "config.ini", "agent config file")
	flag.Parse()
	fmt.Println(*confile)
}
