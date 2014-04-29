package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"syscall"
	"time"
)

type Posts struct {
	XMLName xml.Name `xml:"posts"`
	Posts   []Post   `xml:"post"`
}

type Post struct {
	XMLName xml.Name  `xml:"post"`
	Url     string    `xml:"href,attr"`
	Desc    string    `xml:"description,attr"`
	Notes   string    `xml:"extended,attr"`
	Time    time.Time `xml:"time,attr"`
	Hash    string    `xml:"hash,attr"`
	Shared  bool      `xml:"shared,attr"`
	Tags    string    `xml:"tag,attr"`
	Meta    string    `xml:"meta,attr"`
}

func GetPosts(tags []string) (posts *Posts, err error) {
	options := make(map[string]string)
	options["result"] = "150"
	if len(tags) > 0 {
		options["tag"] = strings.Join(tags, " ")
	}
	url, err := urlWithAuth("/v1/posts/all", options)
	if err != nil {
		return
	}
	res, err := http.Get(url.String())
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()
	xmlString, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	if err = xml.Unmarshal(xmlString, &posts); err != nil {
		return
	}
	return posts, nil
}

func urlWithAuth(pathURL string, options map[string]string) (url.URL, error) {
	u := url.URL{}
	u.Scheme = "https"
	u.Host = path.Join("api.pinboard.in", pathURL)
	q := u.Query()
	if token, ok := syscall.Getenv("PINBOARD_AUTH_TOKEN"); ok {
		q.Set("auth_token", token)
	} else {
		return url.URL{}, fmt.Errorf("`PINBOARD_AUTH_TOKEN` env var is not set.")
	}
	for k, v := range options {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u, nil
}
