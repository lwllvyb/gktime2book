package main

import (
	"flag"

	"github.com/gulywwx/gktime2book/geektime"
)

func main() {
	gk_cellphone := flag.String("cellphone", "1111", "a string")
	gk_password := flag.String("password", "1111", "a string")
	gk_country := flag.String("country", "86", "a string")

	flag.Parse()

	gk := geektime.NewGeekTime(*gk_country, *gk_cellphone, *gk_password)
	// geektime.getIntro(140)
	// geektime.getArticles(140, 100)
	gk.GetArticle(78984)

}
