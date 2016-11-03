package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//TopicCategory struct
type TopicCategory struct {
	CategoryID  string
	Title       string
	Description string
	Topics      []*Topic
}

//InitTopicCategoryList Load All The Category On Start
func InitTopicCategoryList() error {
	TopicsGroupByCategory = TopicsGroupByCategory[:0]
	fp, err := os.OpenFile(categoryJSONFile, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()
	c, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(c, &TopicsGroupByCategory); err != nil {
		return err
	}
	return nil
}
