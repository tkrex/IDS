package controllers

import (
	"github.com/revel/revel"
	"github.com/tkrex/IDS/common/models"
	"net/http"
	"io/ioutil"
)

type Topics struct {
	*revel.Controller
}

var topics = map[int]*models.Topic {
	1: models.NewTopic(1,"test1",[]byte{}),
	2: models.NewTopic(2,"test2",[]byte{}),
}

func (c Topics) List() revel.Result {
	v := make([]*models.Topic, 0, len(topics))

	for  _, value := range topics {
		v = append(v, value)
	}
	return c.RenderJson(v)
}

func (c Topics)  Show(topicID int) revel.Result {

	res,found := topics[topicID]
	if !found {
		return c.NotFound("Could not found Topic")
	}
	return c.RenderJson(res)
}

func (c Topics) Add() revel.Result {

	topic := &models.Topic{}
	if body, err := ioutil.ReadAll(c.Request.Body); err != nil {
		return c.RenderText("bad request")

	} else if err := topic.UnmarshalJSON(body); err != nil {
		return c.RenderText(err.Error())

	}

	existingTopic, ok := topics[topic.ID]
	if ok {
		delete(topics, existingTopic.ID)
	}

	topics[topic.ID] = topic

	c.Response.Status = http.StatusCreated
	return c.RenderJson(topic)
}




