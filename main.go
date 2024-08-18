package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"

	"io"
	"io/ioutil"
	"net/url"
	"path"
	"path/filepath"
	"time"

	"github.com/pengbotao/itopic.go/models"
)

var (
	host           = ""
	isCreateHTML   = false
	isDebug        = false
	isCreateREADME = false
	htmlPrefix     = ""  //without last slash
	htmlDuration   = 300 //5 minutes
	domain         = ""
	githubURL      = "https://github.com/pengbotao/itopic.go"
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1:8001", "host")
	flag.StringVar(&htmlPrefix, "prefix", "../itopic.org", "html folder")
	flag.BoolVar(&isCreateHTML, "html", false, "is create html")
	flag.BoolVar(&isCreateREADME, "readme", false, "is create readme")
	flag.IntVar(&htmlDuration, "duration", 300, "create html duration")
	flag.BoolVar(&isDebug, "debug", false, "debug mode")
}

func main() {
	flag.Parse()
	htmlPrefix = strings.TrimRight(htmlPrefix, "/")
	models.IsDebug = isDebug
	router := loadHTTPRouter()
	ticker := time.NewTicker(1800 * time.Second)
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
		if isDebug == true {
			models.InitTopicList()
			hr := loadHTTPRouter()
			router = hr
		}
		path := r.URL.Path
		if pos := strings.LastIndex(path, "."); pos > 0 {
			path = path[0:pos]
		}
		if bf, ok := router[path]; ok {
			if strings.Compare("/sitemap", path) == 0 {
				w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
			} else {
				w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			}
			w.Write(bf.Bytes())
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		models.InitTopicList()
		hr := loadHTTPRouter()
		router = hr
		http.NotFound(w, r)
	})
	if isCreateREADME == true {
		for i := 0; i < len(models.Topics); i++ {
			fmt.Println("- [" + models.Topics[i].Time.Format("2006-01-02") + "] [" + models.Topics[i].Title + "]" + "(" + url.PathEscape(models.Topics[i].TopicPath) + ")")
		}
		os.Exit(0)
	}
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
	var pages []map[string]string
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
			if t > topicDivCnt {
				isSplit = true
				topicsLeft = models.TopicsGroupByMonth[0:i]
				topicsRight = models.TopicsGroupByMonth[i:]
				break
			}
			t += len(models.TopicsGroupByMonth[i].Topics)
		}
		if isSplit == false {
			topicsLeft = models.TopicsGroupByMonth
		}
	} else {
		topicsLeft = models.TopicsGroupByMonth
	}
	if err := tpl.ExecuteTemplate(&buff, "index.tpl", map[string]interface{}{
		"topics_l":  topicsLeft,
		"topics_r":  topicsRight,
		"domain":    domain,
		"time":      time.Now(),
		"githubURL": githubURL,
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	router["/home"] = buff
	pages = append(pages, map[string]string{
		"loc":        domain + "/",
		"lastmod":    time.Now().Format("2006-01-02"),
		"changefreq": "weekly",
		"priority":   "1",
	})
	//topic router
	for i := range models.Topics {
		var topicLeft = new(models.Topic)
		var topicRight = new(models.Topic)
		if i > 0 {
			topicRight = models.Topics[i-1]
		}
		if i+1 < topicCnt {
			topicLeft = models.Topics[i+1]
		}
		var buff bytes.Buffer
		topicTitle := models.Topics[i].Title
		if isDebug == true {
			re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
			topicTitle = re.ReplaceAllString(topicTitle, "")
		}
		err := tpl.ExecuteTemplate(&buff, "topic.tpl", map[string]interface{}{
			"topic":       models.Topics[i],
			"title":       topicTitle,
			"topic_left":  topicLeft,
			"topic_right": topicRight,
			"domain":      domain,
			"time":        time.Now(),
			"githubURL":   githubURL,
			"i":           i,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		router["/"+models.Topics[i].TopicID] = buff
		pages = append(pages, map[string]string{
			"loc":        domain + "/" + models.Topics[i].TopicID + ".html",
			"lastmod":    models.Topics[i].Time.Format("2006-01-02"),
			"changefreq": "monthly",
			"priority":   "0.9",
		})
	}
	//month router
	for i := range models.TopicsGroupByMonth {
		var buff bytes.Buffer
		err := tpl.ExecuteTemplate(&buff, "list.tpl", map[string]interface{}{
			"title":     models.TopicsGroupByMonth[i].Month,
			"topics":    models.TopicsGroupByMonth[i].Topics,
			"domain":    domain,
			"githubURL": githubURL,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		router["/"+models.TopicsGroupByMonth[i].Month] = buff
		pages = append(pages, map[string]string{
			"loc":        domain + "/" + models.TopicsGroupByMonth[i].Month + ".html",
			"lastmod":    time.Now().Format("2006-01-02"),
			"changefreq": "monthly",
			"priority":   "0.2",
		})
	}
	//tag router
	for i := range models.TopicsGroupByTag {
		var buff bytes.Buffer
		err := tpl.ExecuteTemplate(&buff, "list.tpl", map[string]interface{}{
			"title":     models.TopicsGroupByTag[i].TagName,
			"topics":    models.TopicsGroupByTag[i].Topics,
			"domain":    domain,
			"githubURL": githubURL,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		router["/tag/"+models.TopicsGroupByTag[i].TagID] = buff
		pages = append(pages, map[string]string{
			"loc":        domain + "/tag/" + url.QueryEscape(models.TopicsGroupByTag[i].TagID) + ".html",
			"lastmod":    time.Now().Format("2006-01-02"),
			"changefreq": "monthly",
			"priority":   "0.2",
		})
	}
	//sitemap
	var sitemapBuff bytes.Buffer
	if err := tpl.ExecuteTemplate(&sitemapBuff, "sitemap.tpl", map[string]interface{}{
		"pages": pages,
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	router["/sitemap"] = sitemapBuff
	//create html
	if isCreateHTML == true {
		go generateHTML(router)
		go func() {
			ticker := time.NewTicker(time.Duration(htmlDuration) * time.Second)
			for {
				<-ticker.C
				generateHTML(router)
			}
		}()
	}
	if len(models.Topics) > 0 {
		router["/"] = router["/"+models.Topics[0].TopicID]
		router["/index"] = router["/"]
	}
	return router
}

func generateHTML(router map[string]bytes.Buffer) {
	fmt.Println(time.Now(), "Create Html Start.")
	var err error
	err = clearHTML()
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range router {
		if k == "/" {
			//fmt.Println("Create Html: " + htmlPrefix + k + "index.html")
			writeFile(htmlPrefix+k+"index.html", v)
		} else {
			if k == "/sitemap" {
				//fmt.Println("Create Html: " + htmlPrefix + k + ".xml")
				writeFile(htmlPrefix+k+".xml", v)
			} else {
				//fmt.Println("Create Html: " + htmlPrefix + k + ".html")
				writeFile(htmlPrefix+k+".html", v)
			}
		}
	}
	//copy static folder
	err = copyDir("./static", htmlPrefix+"/static")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Create Html Finished.")
}

func clearHTML() error {
	files, err := ioutil.ReadDir(htmlPrefix)
	if err != nil {
		return err
	}
	for _, file := range files {
		if (file.IsDir() && !strings.HasPrefix(file.Name(), ".")) || strings.HasSuffix(file.Name(), ".html") {
			os.RemoveAll(htmlPrefix + "/" + file.Name())
		}
	}
	return nil
}

func writeFile(filename string, content bytes.Buffer) {
	_, err := os.Stat(path.Dir(filename))
	if os.IsNotExist(err) {
		err := os.MkdirAll(path.Dir(filename), 0775)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	file, _ := os.Create(filename)
	content.WriteTo(file)
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return err
	}

	return
}

func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dst, si.Mode())
		if err != nil {
			return err
		}
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return
}
