package main

import "github.com/tkrex/IDS/daemon/layers"

func main() {
	registrationWoker := layers.NewBrokerRegistrationWorker()
	registrationWoker.GatherBrokerInformation()
}
