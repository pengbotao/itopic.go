package models

import (
	"fmt"
	"os"
)

var (
	topicTagJSONFile    = "posts/tag.json"
	topicMarkdownFolder = "posts"

	//Topics store all the topic
	Topics []*Topic
	//TopicsGroupByMonth store the topic by month
	TopicsGroupByMonth []*MonthList
	//TopicsGroupByTag store all the tag
	TopicsGroupByTag []*TopicTag
)

func init() {
	if err := InitTopicTagList(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := InitTopicList(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
