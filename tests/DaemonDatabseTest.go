package main

import (
	"github.com/tkrex/IDS/daemon/persistence"
	"github.com/tkrex/IDS/common/models"
	"time"
)

func main() {
	dbDelegate,_ := persistence.NewDomainInformationStorage()
	topic := models.NewTopicInformation("testTopic","",time.Now())
	topic.Visibility = true
	topic.Domain = models.NewRealWorldDomain("newDomain")
	dbDelegate.UpdateTopicDomainAndVisibility(topic)
}
