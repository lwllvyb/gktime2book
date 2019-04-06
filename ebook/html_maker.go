package ebook

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/gulywwx/gktime2book/util"
)

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func ParseImage(content string, outDir string) (result string) {
	reg := regexp.MustCompile("img (.{1,15}=\".*?\") src=\".*?\"")
	fucking_styles := reg.FindAllStringSubmatch(content, -1)

	if fucking_styles != nil {
		for _, style := range fucking_styles {
			content = strings.Replace(content, style[1], "", -1)
		}
	}

	reg = regexp.MustCompile("img\\s+src=\"(.*?)\"")
	img_url_list := reg.FindAllStringSubmatch(content, -1)
	if img_url_list != nil {
		for _, url := range img_url_list {
			uuid, _ := newUUID()
			url_local := uuid + ".jpg"

			util.DownloadFile(filepath.Join(outDir, url_local), url[1])
			content = strings.Replace(content, url[1], url_local, -1)
		}
	}

	return content
}

type htmlTemplate struct {
	Title   string
	Content string
}

func renderHTMLFile(title string, content string, outputDir string) {
	data := htmlTemplate{Title: title, Content: content}

	pwd, _ := os.Getwd()

	tmpl, err := template.New("template.html").ParseFiles(filepath.Join(pwd, "ebook", "template.html"))
	if err != nil {
		log.Fatalln(err)
		return
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		log.Fatalln(err)
		return
	}

	out, err := os.Create(filepath.Join(outputDir, title+".html"))
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer out.Close()

	out.Write(tpl.Bytes())

}
