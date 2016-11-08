package models

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"

	"encoding/json"
	"errors"
	"github.com/russross/blackfriday"
)

//InitTopicList load all the topic on init
func InitTopicList() error {
	Topics = Topics[:0]
	TopicsGroupByMonth = TopicsGroupByMonth[:0]
	TopicsGroupByTag = TopicsGroupByTag[:0]
	return filepath.Walk(topicMarkdownFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		t, err := GetTopicByPath(path)
		if err != nil {
			return err
		}
		SetTopicToTag(t)
		SetTopicToMonth(t)
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
	}
	var thj tHeadJSON
	if err := json.Unmarshal([]byte(tHeadStr), &thj); err != nil {
		return nil, errors.New(t.Title + "：" + err.Error())
	}
	t.TopicID = thj.URL
	if t.TopicID == "" {
		return nil, errors.New(t.Title + "：" + err.Error())
	}
	t.Time, err = time.Parse("2006/01/02 15:04", thj.Time)
	if err != nil {
		return nil, errors.New(t.Title + "：" + err.Error())
	}
	if strings.Compare(thj.IsPublic, "no") == 0 {
		t.IsPublic = false
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
	t.Content = string(blackfriday.MarkdownCommon(content.Bytes()))
	return t, nil
}

//SetTopicToTag set topic to tag struct
func SetTopicToTag(t *Topic) {
	if t.IsPublic == false {
		return
	}
	for i := range t.Tag {
		for k := range TopicsGroupByTag {
			if TopicsGroupByTag[k].TagID != t.Tag[i].TagID {
				continue
			}
			isFind := false
			for j := range TopicsGroupByTag[k].Topics {
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
	if t.IsPublic == false {
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
