package ebook

import (
	"log"
	"strings"

	"github.com/gulywwx/gktime2book/geektime"
)

type EBook struct {
	out_dir   string
	course_id int
	gk        *geektime.GeekTime
}

func NewEBook(courseId int, outDir string, gkTime *geektime.GeekTime) *EBook {
	return &EBook{course_id: courseId, out_dir: outDir, gk: gkTime}
}

func (b *EBook) Make() {
	data := b.gk.GetIntro(b.course_id)
	intro := *data
	// log.Println(*intro)
	intro["column_title"] = b.formatTitle(intro["column_title"].(string))
	log.Println(intro["column_title"])
}

func (b *EBook) formatTitle(origin string) (result string) {
	result = strings.Replace(origin, "/", "", -1)
	result = strings.Replace(result, " ", "", -1)
	result = strings.Replace(result, "+", "more", -1)
	result = strings.Replace(result, "\"", "_", -1)
	return result
}
