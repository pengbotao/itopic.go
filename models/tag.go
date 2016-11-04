package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//TopicTag struct
type TopicTag struct {
	TagID       string
	TagName     string
	Description string
	Topics      []*Topic
}

//InitTopicTagList Load All The Tag On Start
func InitTopicTagList() error {
	TopicsGroupByTag = TopicsGroupByTag[:0]
	fp, err := os.OpenFile(topicTagJSONFile, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()
	c, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(c, &TopicsGroupByTag); err != nil {
		return err
	}
	return nil
}
