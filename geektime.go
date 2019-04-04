package main

import "net/http"

type geektime struct {
	client    *http.Client
	cookie    string
	country   string
	cellphone string
	password  string
}

func NewGeekTime(country string, cellphone string, password string) *geektime {
	return &geektime{
		country:   country,
		cellphone: cellphone,
		password:  password,
	}
}

func getCookie() {}
