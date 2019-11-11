package ebook

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/lwllvyb/gktime2book/util"

	"github.com/mattn/godown"
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

func ParseImage(content string, outDir string, relative string) (result string) {
	reg := regexp.MustCompile("img (.{1,15}=\".*?\") src=\".*?\"")
	fuckingStyles := reg.FindAllStringSubmatch(content, -1)

	if fuckingStyles != nil {
		for _, style := range fuckingStyles {
			content = strings.Replace(content, style[1], "", -1)
		}
	}

	reg = regexp.MustCompile("img\\s+src=\"(.*?)\"")
	imgURLList := reg.FindAllStringSubmatch(content, -1)
	if imgURLList != nil {
		for _, url := range imgURLList {
			uuid, _ := newUUID()
			urlLocal := uuid + ".jpg"
			if err := util.DownloadFile(filepath.Join(outDir, urlLocal), url[1]); err != nil {
				log.Printf("DownloadFile [%s] -> %s error: %v\n", url[1], filepath.Join(outDir, urlLocal))
			}

			content = strings.Replace(content, url[1], filepath.Join(relative, urlLocal), -1)
		}
	}

	return content
}

func ParseAudio(audioHTTP, outDir string) string {
	u, err := url.Parse(audioHTTP)
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	audioName := path.Base(u.Path)
	downloadPath := filepath.Join(outDir, audioName)
	if err := util.DownloadFile(downloadPath, audioHTTP); err != nil {
		log.Printf("DownloadFile [%s] -> %s error:%v", downloadPath, audioHTTP, err)
	}
	return audioName
}

type htmlTemplate struct {
	Title   string
	Content string
}

func renderHTML(title string, content string) (bytes.Buffer, error) {
	data := htmlTemplate{Title: title, Content: content}

	pwd, _ := os.Getwd()

	tmpl, err := template.New("template.html").ParseFiles(filepath.Join(pwd, "ebook", "template.html"))
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("New error: %v", err)
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return bytes.Buffer{}, fmt.Errorf("Execute error: %v", err)
	}
	return tpl, nil
}
func writeHTMLFile(dstFile, title string, content []byte) {

	out, err := os.Create(dstFile)
	if err != nil {
		log.Fatalf("create %v\n", err)
		return
	}
	defer out.Close()

	out.Write(content)

}

func genHTMLFile(dstFile, title, content string) {
	tpl, err := renderHTML(title, content)
	if err != nil {
		log.Fatalf("renderHTML error:%v", err)
		return
	}
	writeHTMLFile(dstFile, title, tpl.Bytes())
}

func renderH(title, content string) string {
	return fmt.Sprintf("<h1>%s</h1>\n%s", title, content)
}

func genMardown(dstFile, title, content, audioFile string) {
	// tpl, err := renderHTML(title, content)
	// if err != nil {
	// 	log.Fatalf("renderHTML error:%v", err)
	// 	return
	// }
	f, err := os.OpenFile(dstFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	headerData := fmt.Sprintf("# %s\n", title)
	headerData += genAudioData(audioFile)

	if _, err = f.WriteString(headerData); err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	rBuf := bytes.NewReader([]byte(content))
	godown.Convert(&buf, rBuf, nil)
	if _, err = f.WriteString(buf.String()); err != nil {
		panic(err)
	}
}

func genBookJSONFile(dstFile, columnTitle, language string) {
	data := map[string]string{
		"title":    columnTitle,
		"language": language,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("json.Marshal fail, error:%v", err)
		return
	}
	ioutil.WriteFile(dstFile, jsonData, os.ModePerm)
}

func genSummaryData(name, dstFile string) string {
	return fmt.Sprintf("* [%s](%s)\n", name, dstFile)
}

func genAudioData(audioFile string) string {
	if len(audioFile) == 0 {
		return ""
	}
	return fmt.Sprintf(`<audio id="audio" controls="" preload="none"><source id="mp3" src="%s"></audio>`, audioFile) + "\n\n"
}
