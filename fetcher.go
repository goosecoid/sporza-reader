package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iunary/fakeuseragent"
)

func fetchPageFromUrl(url string) (*http.Response, error) {
	agent := fakeuseragent.GetUserAgent(fakeuseragent.BrowserFirefox)
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", agent)
	res, err := client.Do(req)

	if res.StatusCode != 200 {
		return nil, err
	}

	return res, nil
}

func fetchUrls(url string) ([]Article, error) {
	var articles []Article
	count := 0
	resp, err := fetchPageFromUrl(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	filters := getFilters()
	doc.Find("main#content").Each(func(i int, s *goquery.Selection) {
		var url string

		s.Find("a").Each(func(i int, s *goquery.Selection) {
			val, _ := s.Attr("class")

			if strings.Contains(val, "_card") {
				url = s.AttrOr("href", "")

				s.Find("div").Each(func(i int, s *goquery.Selection) {
					val, _ := s.Attr("class")

					if strings.Contains(val, "_title") {

						title := s.Text()
						lowerTitle := strings.ToLower(title)
						skip := false

						for _, filter := range filters {
							if strings.Contains(lowerTitle, filter) {
								skip = true
							}
						}

						if !skip {
							formattedUrl := fmt.Sprintf("http://localhost:8080/article?id=%d", count)
							articles = append(articles, Article{Url: url, FormattedUrl: formattedUrl, Title: title, ID: count, Content: ""})
							count++
						}
					}
				})
			}
		})
	})

	if err != nil {
		return nil, err
	}

	return articles, nil
}

func fetchSanitizedHtmlFromUrl(url string) (string, error) {

	res, err := fetchPageFromUrl(url)

	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	defer res.Body.Close()

	if err != nil {
		return "", err
	}

	var content string

	doc.Find("section").Each(func(i int, s *goquery.Selection) {
		if s.Children().First().Text() == "Gerelateerd:" {
			s.Remove()
		}
	})

	doc.Find("span:contains('meer tonen')").Remove()

	doc.Find("div[aria-live='polite']").Remove()

	doc.Find("h2:contains('Beluister de analyse')").Remove()

	doc.Find("h2:contains('Fase per fase')").Parent().Remove()

	doc.Find(".sw-article-layout-main").Each(func(i int, s *goquery.Selection) {
		html, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}
		content = html
	})

	p := SporzaPolicy()
	sanitizedHtml := p.Sanitize(content)

	return sanitizedHtml, nil
}
