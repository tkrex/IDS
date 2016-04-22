package models

import (
	"time"
)

type Topic struct {
	ID		    int        `json:"id"`
	Name                string        `json:"name"`
	LastPayload         []byte        `json:"payload"`
	LastUpdateTimeStamp time.Time        `json:"lastUpdate"`
	UpdateInterval      int      `json:"updateInterval"`
	NumberOfUpdates	    int		  `json:"numberOfUpdates"`
	Domain		RealWorldDomain   `json:"domain"`
}



func NewTopic(id int, name string, payload []byte) *Topic {
	topic := new(Topic)
	topic.ID = id
	topic.Name = name
	topic.LastPayload = payload
	topic.LastUpdateTimeStamp = time.Now()
	return topic
}

//func (t *Topic) UnmarshalJSON(data []byte) error {
//	if (t == nil) {
//		return errors.New("Structure: UnmarshalJSON on nil pointer")
//	}
//	var fields map[string]interface{}
//	json.Unmarshal(data, &fields)
//	id := int(fields["id"].(float64))
//
//	name, errName := fields["name"].(string)
//
//	if !errName  {
//		return errors.New("Name Parsing error")
//	}
//	//t.LastPayload= fields["payload"].([]byte)
//	//t.LastUpdateTimeStamp = time.fields["lastUpdate"].(string)
//	//t.UpdateInterval = fields["updateInterval"].(int)
//	//t.NumberOfUpdates = fields["numberOfUpdates"].(int)
//	t.Name = name
//	t.ID = id
//	var domain RealWorldDomain
//	if err := json.Unmarshal(fields["domain"].([]byte),&domain); err != nil {
//		return err
//	}
//	t.Domain = domain
//	return nil
//}

