package ebook

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gulywwx/gktime2book/geektime"
	"github.com/gulywwx/gktime2book/util"
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
	column_title := b.formatTitle(intro["column_title"].(string))
	column_intro := intro["column_intro"].(string)
	column_cover := intro["column_cover"].(string)

	list := b.gk.GetArticles(b.course_id, 1000)
	articles := *list

	output_dir := b.out_dir + string(os.PathSeparator) + column_title
	if _, err := os.Stat(output_dir); err == nil {
		os.RemoveAll(output_dir)
	}
	_ = os.Mkdir(output_dir, os.ModePerm)

	renderHTMLFile("简介", ParseImage(column_intro, output_dir), output_dir)
	log.Println("下载", column_title, "简介", "done")

	util.DownloadFile(filepath.Join(output_dir, "cover.jpg"), column_cover)
	log.Println("下载", column_title, "封面", "done")

	for _, a := range articles {
		var article map[string]interface{} = a.(map[string]interface{})
		article_title := b.formatTitle(article["article_title"].(string))
		id := article["id"].(float64)

		data := b.gk.GetArticle(int(id))
		article = *data
		article_content := article["article_content"].(string)
		// log.Println(article_title, id, article_content)
		renderHTMLFile(article_title, ParseImage(article_content, output_dir), output_dir)
		log.Println("下载", column_title, ":", article_title, "done")
	}

}

func (b *EBook) formatTitle(origin string) (result string) {
	result = strings.Replace(origin, "/", "", -1)
	result = strings.Replace(result, " ", "", -1)
	result = strings.Replace(result, "+", "more", -1)
	result = strings.Replace(result, "\"", "_", -1)
	result = strings.Replace(result, "|", "_", -1)

	return result
}
