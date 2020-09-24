package models

import (
	"fmt"
	"os"
	"time"
)

//Topic struct
type Topic struct {
	TopicID        string
	Title          string
	Time           time.Time
	LastModifyTime time.Time
	Tag            []*TopicTag
	Content        string
	TopicPath      string
	IsPublic       bool //true for public，false for protected
	IsToc          bool
}

//TopicTag struct
type TopicTag struct {
	TagID   string
	TagName string
	Topics  []*Topic
}

//TopicMonth show the topic group by month
type TopicMonth struct {
	Month  string
	Topics []*Topic
}

var (
	//IsDebug assign from main.go
	IsDebug             = false
	topicMarkdownFolder = "posts"
	randomWrittingColor = []string{"#FE9A2E", "#BFFF00", "#81F7BE", "#2ECCFA", "#F781F3", "#A9A9F5", "#D8D8D8", "#FF4000", "#F4FA58", "#5858FA"}
	//Topics store all the topic
	Topics []*Topic
	//TopicsGroupByMonth store the topic by month
	TopicsGroupByMonth []*TopicMonth
	//TopicsGroupByTag store all the tag
	TopicsGroupByTag []*TopicTag
)

func init() {
	if err := InitTopicList(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
