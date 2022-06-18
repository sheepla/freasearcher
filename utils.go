package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skratchdot/open-golang/open"
	"github.com/tidwall/gjson"
)

type Result struct {
	title   string
	content string
	url     string
}

type Param struct {
	Query      string
	Language   string
	SafeSearch int
}

func newURL(param Param) string {
	u := &url.URL{
		Scheme: "https",
		Host:   "freasearch.org",
		Path:   "search",
	}
	q := u.Query()
	q.Set("format", "json")
	q.Set("q", param.Query)
	q.Set("language", param.Language)
	q.Set("safesearch", strconv.Itoa(param.SafeSearch))
	u.RawQuery = q.Encode()

	return u.String()
}

func getResp(param Param) ([]Result, error) {
	req, err := http.NewRequest("GET", newURL(param), nil)
	if err != nil {
		return nil, fmt.Errorf("リクエストの構築に失敗しました: %w", err)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエストの取得に失敗しました: %w", err)
	}

	defer resp.Body.Close()

	bArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスボディの読み込みに失敗しました: %w", err)
	}

	results := gjson.Get(string(bArray), "results")

	ctns := []Result{}

	for _, result := range results.Array() {

		title := gjson.Get(result.String(), "title").String()
		content := gjson.Get(result.String(), "content").String()
		url := gjson.Get(result.String(), "url").String()

		tmp := Result{title: title, content: content, url: url}
		ctns = append(ctns, tmp)

	}

	return ctns, nil
}

func openBrowser(url string) {
	fmt.Printf("Open the %s in your browser...", url)
	open.Run(url)
}
