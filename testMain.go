package main

import (
)
import "github.com/tkrex/IDS/daemon/layers"

func main() {
	_ = layers.NewBrokerRegistrationWorker()
	for {}
}
