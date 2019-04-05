package geektime

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const host = "https://time.geekbang.org/serv/v1"

type GeekTime struct {
	client    *http.Client
	cookie    string
	country   string
	cellphone string
	password  string
	links     map[string]string
}

func NewGeekTime(country string, cellphone string, password string) *GeekTime {
	return &GeekTime{
		country:   country,
		cellphone: cellphone,
		password:  password,
		client: &http.Client{Timeout: time.Second * 1000, Transport: &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}},
		links: map[string]string{
			"login":       "https://account.geekbang.org/account/ticket/login",
			"products":    host + "/my/products/all",
			"productList": host + "/my/products/list",
			"intro":       host + "/column/intro",
			"articles":    host + "/column/articles",
			"article":     host + "/article",
			"comments":    host + "/comments",
			"audios":      host + "/column/audios"},
	}
}

func (g *GeekTime) request(url string, payload *map[string]interface{}, cookie string) (data interface{}, loginCookie string) {

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(g.payload(payload)))
	if err != nil {
		log.Fatalln(err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Referer", url)
	request.Header.Set("Cookie", cookie)

	response, err := g.client.Do(request)
	if err != nil {
		log.Fatalln(err)
		return
	}

	loginCookie = ""
	if url == g.links["login"] {
		// for k, v := range response.Header {
		// 	log.Println("key=", k, "value=", v[0])
		// 	if k == "Set-Cookie" {
		// 		if loginCookie == "" {
		// 			loginCookie = v[0]
		// 		}
		// 		loginCookie = loginCookie + "; " + v[0]
		// 	}

		// }

		for _, cookie := range response.Cookies() {
			if loginCookie == "" {
				loginCookie = cookie.Name + "=" + cookie.Value
			}
			loginCookie = loginCookie + "; " + cookie.Name + "=" + cookie.Value
		}
	}

	body := response.Body
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatalln(err)
	}

	var temp map[string]interface{}
	err = json.Unmarshal(b, &temp)
	if err != nil {
		log.Fatalln(err)
	}

	return temp["data"], loginCookie
}

func (g *GeekTime) payload(payload *map[string]interface{}) (body []byte) {
	bytesRepresentation, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return bytesRepresentation
}

func (g *GeekTime) getCookie() (cookie string) {
	if g.cookie != "" {
		return g.cookie
	}

	var payload = map[string]interface{}{
		"country":   g.country,
		"cellphone": g.cellphone,
		"password":  g.password,
		"captcha":   "",
		"remember":  1,
		"platform":  3,
		"appid":     1,
	}

	_, cookie = g.request(g.links["login"], &payload, "")

	g.cookie = cookie
	return g.cookie
}

func (g *GeekTime) GetIntro(cid int) interface{} {
	cookie := g.getCookie()
	var payload = map[string]interface{}{
		"cid":           cid,
		"with_groupbuy": true,
	}

	data, _ := g.request(g.links["intro"], &payload, cookie)

	return data

}

func (g *GeekTime) GetArticles(cid int, size int) interface{} {
	cookie := g.getCookie()
	var payload = map[string]interface{}{
		"cid":    cid,
		"size":   size,
		"order":  "earliest",
		"prev":   0,
		"sample": false,
	}

	data, _ := g.request(g.links["articles"], &payload, cookie)
	return data
}

func (g *GeekTime) GetArticle(id int) interface{} {
	cookie := g.getCookie()
	var payload = map[string]interface{}{
		"id":                id,
		"include_neighbors": false,
	}
	data, _ := g.request(g.links["article"], &payload, cookie)

	a, _ := data.(map[string]interface{})
	log.Println(a["article_title"])
	// log.Println(a["article_content"])
	return data
}
