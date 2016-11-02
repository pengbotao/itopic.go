package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type TopicCategory struct {
	CategoryId  string
	Title       string
	Description string
	Topics      []*Topic
}

var Categories []*TopicCategory

func InitTopicCategoryList() error {
	Categories = Categories[:0]
	fp, err := os.OpenFile("posts/category.json", os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()
	c, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(c, &Categories); err != nil {
		return err
	}
	return nil
}
