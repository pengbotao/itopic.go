package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"itopic.go/models"
)

var (
	host   = "127.0.0.1:8001"
	router = make(map[string]bytes.Buffer)
)

func init() {
	if err := models.InitTopicCategoryList(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := models.InitTopicList(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var tpl *template.Template
	tpl, err := template.ParseGlob("views/*.tpl")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//topic router
	for i := range models.Topics {
		if models.Topics[i].IsPublic == false {
			continue
		}
		var buff bytes.Buffer
		err := tpl.ExecuteTemplate(&buff, "topic.tpl", map[string]interface{}{
			"topic": models.Topics[i],
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		router["/"+models.Topics[i].TopicID] = buff
	}
	//category router
	for i := range models.Categories {
		var buff bytes.Buffer
		err := tpl.ExecuteTemplate(&buff, "list.tpl", map[string]interface{}{
			"title":  models.Categories[i].Title,
			"topics": models.Categories[i].Topics,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		router["/category/"+models.Categories[i].CategoryID] = buff
	}
	//month router
	for i := range models.TopicGroupByMonth {
		var buff bytes.Buffer
		err := tpl.ExecuteTemplate(&buff, "list.tpl", map[string]interface{}{
			"title":  models.TopicGroupByMonth[i].Month,
			"topics": models.TopicGroupByMonth[i].Topics,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		router["/"+models.TopicGroupByMonth[i].Month] = buff
	}
	//homepage router
	var buff bytes.Buffer
	if err := tpl.ExecuteTemplate(&buff, "index.tpl", map[string]interface{}{
		"topics_m": models.TopicGroupByMonth,
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	router["/"] = buff
}

func main() {
	http.Handle("/favicon.ico", http.FileServer(http.Dir("public")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("public/uploads"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		path := r.URL.Path
		if pos := strings.LastIndex(path, "."); pos > 0 {
			path = path[0:pos]
		}
		if bf, ok := router[path]; ok {
			w.Write(bf.Bytes())
		} else {
			http.NotFound(w, r)
		}
	})

	fmt.Println("The topic server is running at http://" + host)
	fmt.Println("Quit the server with Control-C")
	http.ListenAndServe(host, nil)
}
