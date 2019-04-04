package main

import (
	"flag"
	"fmt"
)

func main() {
	gk_cellphone := flag.String("cellphone", "1111", "a string")
	gk_password := flag.String("password", "1111", "a string")
	gk_country := flag.String("country", "86", "a string")

	flag.Parse()

	fmt.Println("gk_cellphone:", *gk_cellphone)
	fmt.Println("gk_password:", *gk_password)
	fmt.Println("gk_country:", *gk_country)

}
