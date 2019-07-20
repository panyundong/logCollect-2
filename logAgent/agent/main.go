package main

import (
	"fmt"
	_ "github.com/astaxie/beego"
	"rsc.io/quote"
)

func main() {
	fmt.Println(quote.Hello())
}
