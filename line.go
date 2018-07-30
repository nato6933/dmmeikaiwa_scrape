package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	notify_url = "https://notify-api.line.me/api/notify"
)

type line struct {
	accessToken string
	Url         string
}

func newLine(token string) *line {
	return &line{
		accessToken: token,
		Url:         notify_url,
	}
}

func (lobj *line) notify(msg string) {

	if lobj.accessToken == "" {
		log.Printf("accessToken required, but not set.")
		return
	}

	if msg == "" {
		log.Printf("msg required, but not set.")
		return
	}

	u, err := url.ParseRequestURI(lobj.Url)
	if err != nil {
		log.Fatal(err)
	}

	c := &http.Client{}

	form := url.Values{}
	form.Add("message", msg)

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+lobj.accessToken)

	_, err = c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

}
