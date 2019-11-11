package ebook

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gktime2book/geektime"
)

type EBook struct {
	outDir   string
	courseID int
	gk       *geektime.GeekTime
}

func NewEBook(courseId int, outDir string, gkTime *geektime.GeekTime) *EBook {
	return &EBook{courseID: courseId, outDir: outDir, gk: gkTime}
}

func (b *EBook) Make() {
	// d := b.gk.GetAllColumns()
	// log.Println(d)
	// return
	data := b.gk.GetIntro(b.courseID)
	intro := *data
	columnTitle := b.formatTitle(intro["column_title"].(string))
	columnIntro := intro["column_intro"].(string)
	// column_cover := intro["column_cover"].(string)

	list := b.gk.GetArticles(b.courseID, 1000)
	articles := *list
	columeOutDir := filepath.Join(b.outDir, columnTitle)

	docDir := filepath.Join(columeOutDir, "docs")

	_ = os.MkdirAll(docDir, os.ModePerm)

	summary := "# SUMMARY\n\n"
	summary += genSummaryData(columnTitle, "README.md")

	// README.md
	readmeFile := filepath.Join(columeOutDir, "README.md")
	genMardown(readmeFile, columnTitle, ParseImage(columnIntro, docDir, "docs"), "")
	// book.json
	bookjsonFile := filepath.Join(columeOutDir, "book.json")
	genBookJSONFile(bookjsonFile, columnTitle, "zh")

	log.Println("下载", columnTitle, columnTitle, "done")

	for _, a := range articles {
		article := a.(map[string]interface{})
		articleTitle := article["article_title"].(string)
		id := article["id"].(float64)

		data := b.gk.GetArticle(int(id))
		article = *data
		// log.Println(article)
		if article["article_content"] == nil {
			break
		}
		articleMD := strconv.FormatInt(int64(id), 10) + ".md"
		articleMDPath := filepath.Join(docDir, articleMD)
		if _, err := os.Stat(articleMDPath); err == nil {
			log.Printf("%s:%s already done\n", columnTitle, articleTitle)
		} else {
			articleContent := article["article_content"].(string)
			content := ParseImage(articleContent, docDir, "")
			audioHTTP := article["audio_download_url"].(string)
			audioName := ParseAudio(audioHTTP, docDir)
			genMardown(articleMDPath, articleTitle, content, audioName)
			log.Printf("Download %s:%s done\n", columnTitle, articleTitle)
		}
		summary += genSummaryData(articleTitle, "./docs/"+articleMD)
	}
	// SUMMARY.md
	ioutil.WriteFile(filepath.Join(columeOutDir, "SUMMARY.md"), []byte(summary), os.ModePerm)
}

func (b *EBook) formatTitle(origin string) (result string) {
	result = strings.Replace(origin, "/", "", -1)
	result = strings.Replace(result, " ", "", -1)
	result = strings.Replace(result, "+", "more", -1)
	result = strings.Replace(result, "\"", "_", -1)
	result = strings.Replace(result, "|", "_", -1)

	return result
}
