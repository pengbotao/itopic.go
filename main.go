package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"itopic.go/models"
	"time"
)

var (
	host = "127.0.0.1:8001"
)

func main() {
	router := loadHTTPRouter()
	ticker := time.NewTicker(3600 * time.Second)
	go func() {
		for range ticker.C {
			models.InitTopicList()
			hr := loadHTTPRouter()
			router = hr
		}
	}()
	http.Handle("/favicon.ico", http.FileServer(http.Dir("static")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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

	fmt.Printf("The topic server is running at http://%s\n", host)
	fmt.Printf("Quit the server with Control-C\n\n")
	if err := http.ListenAndServe(host, nil); err != nil {
		fmt.Print(err)
	}
}

func loadHTTPRouter() map[string]bytes.Buffer {
	router := make(map[string]bytes.Buffer)
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
	//tag router
	for i := range models.TopicsGroupByTag {
		var buff bytes.Buffer
		err := tpl.ExecuteTemplate(&buff, "list.tpl", map[string]interface{}{
			"title":  models.TopicsGroupByTag[i].TagName,
			"topics": models.TopicsGroupByTag[i].Topics,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		router["/tag/"+models.TopicsGroupByTag[i].TagID] = buff
	}
	//month router
	for i := range models.TopicsGroupByMonth {
		var buff bytes.Buffer
		err := tpl.ExecuteTemplate(&buff, "list.tpl", map[string]interface{}{
			"title":  models.TopicsGroupByMonth[i].Month,
			"topics": models.TopicsGroupByMonth[i].Topics,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		router["/"+models.TopicsGroupByMonth[i].Month] = buff
	}
	//homepage router
	var buff bytes.Buffer
	topicCnt := len(models.Topics)
	topicDivCnt := topicCnt / 2
	var topicsLeft []*models.TopicMonth
	var topicsRight []*models.TopicMonth
	if topicDivCnt > 0 {
		t := 0
		isSplit := false
		for i := range models.TopicsGroupByMonth {
			t += len(models.TopicsGroupByMonth[i].Topics)
			if t > topicDivCnt {
				isSplit = true
				topicsLeft = models.TopicsGroupByMonth[0:i]
				topicsRight = models.TopicsGroupByMonth[i:]
				break
			}
		}
		if isSplit == false {
			topicsLeft = models.TopicsGroupByMonth
		}
	} else {
		topicsLeft = models.TopicsGroupByMonth
	}
	if err := tpl.ExecuteTemplate(&buff, "index.tpl", map[string]interface{}{
		"topics_l": topicsLeft,
		"topics_r": topicsRight,
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	router["/"] = buff
	return router
}
