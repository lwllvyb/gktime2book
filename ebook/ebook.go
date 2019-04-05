package ebook

import (
	"log"

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
	intro := b.gk.GetIntro(b.course_id)
	log.Println(intro)

}
