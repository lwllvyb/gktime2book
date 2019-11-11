package main

import (
	"flag"

	"gktime2book/ebook"
	"gktime2book/geektime"
)

//go run main.go --cellphone=xxxxxxx --password=*****
func main() {
	gk_cellphone := flag.String("cellphone", "0", "a string")
	gk_password := flag.String("password", ",0", "a string")
	gk_country := flag.String("country", "86", "a string")
	gk_cid := flag.Int("cid", 0, "coursel id")

	flag.Parse()

	gk := geektime.NewGeekTime(*gk_country, *gk_cellphone, *gk_password)
	if *gk_cid == 0 {
		allColumns := gk.GetAllColumns()
		columns := allColumns["1"]
		for _, id := range columns {
			// eb := ebook.NewEBook(140, "d:", gk)
			eb := ebook.NewEBook(id, "./gitbook", gk)
			eb.Make()

		}
	} else {
		eb := ebook.NewEBook(*gk_cid, "./gitbook", gk)
		eb.Make()
	}

	// geektime.GetIntro(140)
	// geektime.GetArticles(140, 100)
	// data := gk.GetArticle(87980)
	// log.Println(*data)

}
