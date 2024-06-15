package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Article struct {
	ID           int    `json:"id"`
	Url          string `json:"url"`
	FormattedUrl string `json:"formattedUrl"`
	Title        string `json:"title"`

	// Content should rather be of type template.HTML
	Content template.HTML `json:"content"`
}

var articlesMap map[int]Article = make(map[int]Article)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		articles, err := fetchUrls("https://sporza.be/nl/pas-verschenen")

		if err != nil {
			log.Fatal(err)
		}

		for i, v := range articles {
			articlesMap[v.ID] = articles[i]
		}

		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, articles)
	})

	http.HandleFunc("/article", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		i, err := strconv.Atoi(id)
		if err != nil {
			log.Fatal(err)
		}

		article := articlesMap[i]

		content, err := fetchSanitizedHtmlFromUrl(article.Url)

		if err != nil {
			log.Fatal(err)
		}

		t, _ := template.ParseFiles("templates/article.html")
		t.Execute(w, Article{ID: article.ID, Url: article.Url, FormattedUrl: article.FormattedUrl, Title: article.Title, Content: template.HTML(content)})
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
