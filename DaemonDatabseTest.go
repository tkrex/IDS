package main

import (
	"github.com/tkrex/IDS/daemon/persistence"
	"github.com/tkrex/IDS/common/models"
	"time"
)

func main() {
	dbDelegate,_ := persistence.NewDaemonDatabaseWorker()
	topic := models.NewTopic("testTopic","",time.Now())
	topic.Visibility = true
	topic.Domain = models.NewRealWorldDomain("newDomain")
	dbDelegate.UpdateTopicDomainAndVisibility(topic)
}
