package models

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"encoding/json"
	"errors"

	"github.com/russross/blackfriday"
)

const (
	markdownHTMLFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	markdownExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_AUTO_HEADER_IDS
)

//InitTopicList load all the topic on init
func InitTopicList() error {
	topicsMutex.Lock()
	defer topicsMutex.Unlock()

	Topics = make([]*Topic, 0)        // 不保留容量
	TopicsGroupByMonth = make([]*TopicMonth, 0)
	TopicsGroupByTag = make([]*TopicTag, 0)

	fileCount := 0
	successCount := 0
	err := filepath.Walk(topicMarkdownFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		fileCount++
		t, err := GetTopicByPath(path)
		if err != nil {
			if IsDebug {
				fmt.Printf("Error loading %s: %v\n", path, err)
			}
			return nil
		}
		if t.TopicPath == "" {
			if IsDebug {
				fmt.Printf("Skipping %s: empty TopicPath\n", path)
			}
			return nil
		}
		if IsDebug == false && t.IsPublic == false {
			return nil
		}
		successCount++
		SetTopicToTag(t)
		SetTopicToMonth(t)
		
		// 添加nil检查
		if t == nil || t.Time.IsZero() {
			return nil
		}
		
		//append topics desc
		for i := range Topics {
			if t.Time.After(Topics[i].Time) {
				Topics = append(Topics, nil)
				copy(Topics[i+1:], Topics[i:])
				Topics[i] = t
				return nil
			}
		}
		Topics = append(Topics, t)
		return nil
	})
	if IsDebug {
		fmt.Printf("InitTopicList: scanned %d files, loaded %d topics, %d private skipped\n",
			fileCount, successCount, fileCount-successCount)
	}
	return err
}

//GetTopicByPath read the topic by path
func GetTopicByPath(path string) (*Topic, error) {
	fp, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.New(path + "：" + err.Error())
	}
	defer fp.Close()
	t := &Topic{
		Title:    strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		IsPublic: true,
	}
	var tHeadStr string
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		s := scanner.Text()
		tHeadStr += s
		if len(s) == 0 {
			break
		}
	}
	tHeadStr = strings.Trim(tHeadStr, "```")
	type tHeadJSON struct {
		URL      string
		Time     string
		Tag      string
		IsPublic string `json:"public"`
		IsToc    string `json:"toc"`
	}
	var thj tHeadJSON
	if err := json.Unmarshal([]byte(tHeadStr), &thj); err != nil {
		fmt.Println("Notice: " + path + "/" + t.Title + "：" + err.Error())
		return t, nil
	}
	t.TopicID = thj.URL
	if t.TopicID == "" {
		return nil, errors.New(t.Title + "：" + err.Error())
	}
	t.Time, err = time.Parse("2006/01/02 15:04", thj.Time)
	if err != nil {
		// 如果解析失败，使用文件修改时间
		if info, err := os.Stat(path); err == nil {
			t.Time = info.ModTime()
		} else {
			t.Time = time.Now() // 最后备选
		}
	}
	if strings.Compare(thj.IsPublic, "no") == 0 {
		t.IsPublic = false

		t.Title = "<font color=\"" + randomWritingColor[rand.Intn(len(randomWritingColor))] + "\">" + t.Title + "</font>"
	}
	tagArray := strings.Split(thj.Tag, ",")
	var isFind bool
	for _, tagName := range tagArray {
		tagName = strings.TrimSpace(tagName)
		if len(tagName) == 0 {
			continue
		}
		isFind = false
		for kc := range TopicsGroupByTag {
			if strings.Compare(strings.ToLower(tagName), TopicsGroupByTag[kc].TagID) == 0 {
				t.Tag = append(t.Tag, TopicsGroupByTag[kc])
				isFind = true
				break
			}
		}
		if isFind == false {
			tt := &TopicTag{TagID: strings.ToLower(tagName), TagName: tagName}
			t.Tag = append(t.Tag, tt)
			TopicsGroupByTag = append(TopicsGroupByTag, tt)
		}
	}
	var content bytes.Buffer
	for scanner.Scan() {
		content.Write(scanner.Bytes())
		content.WriteString("\n")
	}
	var mdFlag = markdownHTMLFlags
	t.IsToc = false
	if thj.IsToc == "yes" {
		t.IsToc = true
		mdFlag = mdFlag | blackfriday.HTML_TOC
	}
	var markdownRenderer = blackfriday.HtmlRenderer(mdFlag, "", "")
	t.Content = string(blackfriday.MarkdownOptions(content.Bytes(), markdownRenderer, blackfriday.Options{Extensions: markdownExtensions}))
	t.TopicPath = path

	finfo, _ := os.Stat(path)
	lastModTime := finfo.ModTime()
	if lastModTime.Unix()-t.Time.Unix() > 7*86400 && time.Now().Unix()-lastModTime.Unix() < 365*86400 {
		t.LastModifyTime = lastModTime
	}
	return t, nil
}

//SetTopicToTag set topic to tag struct
func SetTopicToTag(t *Topic) {
	if IsDebug == false && t.IsPublic == false {
		return
	}
	if t == nil {
		return
	}
	for i := range t.Tag {
		if t.Tag[i] == nil {
			continue
		}
		
		for k := range TopicsGroupByTag {
			if TopicsGroupByTag[k] == nil || TopicsGroupByTag[k].TagID != t.Tag[i].TagID {
				continue
			}
			isFind := false
			for j := range TopicsGroupByTag[k].Topics {
				if TopicsGroupByTag[k].Topics[j] == nil {
					continue
				}
				if t.Time.After(TopicsGroupByTag[k].Topics[j].Time) {
					TopicsGroupByTag[k].Topics = append(TopicsGroupByTag[k].Topics, nil)
					copy(TopicsGroupByTag[k].Topics[j+1:], TopicsGroupByTag[k].Topics[j:])
					TopicsGroupByTag[k].Topics[j] = t
					isFind = true
					break
				}
			}
			if isFind == false {
				TopicsGroupByTag[k].Topics = append(TopicsGroupByTag[k].Topics, t)
			}
			break
		}
	}
}

//SetTopicToMonth set topic to month struct
func SetTopicToMonth(t *Topic) {
	if IsDebug == false && t.IsPublic == false {
		return
	}
	month := t.Time.Format("2006-01")
	tm := &TopicMonth{}
	for _, m := range TopicsGroupByMonth {
		if m.Month == month {
			tm = m
		}
	}
	if tm.Month == "" {
		tm.Month = month
		isFind := false
		for i := range TopicsGroupByMonth {
			if strings.Compare(tm.Month, TopicsGroupByMonth[i].Month) > 0 {
				TopicsGroupByMonth = append(TopicsGroupByMonth, nil)
				copy(TopicsGroupByMonth[i+1:], TopicsGroupByMonth[i:])
				TopicsGroupByMonth[i] = tm
				isFind = true
				break
			}
		}
		if isFind == false {
			TopicsGroupByMonth = append(TopicsGroupByMonth, tm)
		}
	}
	for i := range tm.Topics {
		if t.Time.After(tm.Topics[i].Time) {
			tm.Topics = append(tm.Topics, nil)
			copy(tm.Topics[i+1:], tm.Topics[i:])
			tm.Topics[i] = t
			return
		}
	}
	tm.Topics = append(tm.Topics, t)
}
