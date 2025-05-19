package utils

import (
	"sim_data_gen/messaging/nats"
	"fmt"
)

var natPublisher nats.Publisher
var natSub nats.PubSub

func GetNatPublisher() nats.Publisher {
	return natPublisher
}
func GetNatSubPub() nats.PubSub {
	return natSub
}
func DisConnectNats() {
	if natPublisher != nil {
		natPublisher.Close()
		natPublisher = nil
	}
	if natSub != nil {
		natSub.Close()
		natSub = nil
	}
}
func ConnectNats(natsVarName string) {
	var err error
	var natsUrl string
	if natsUrl = GetConfig().String(natsVarName+".url", ""); natsUrl == "" {
		Logger.Error("nats url is empty")
		return
	}

	if natPublisher, err = nats.NewPublisher(natsUrl); err != nil {
		Logger.Error(err.Error())
		return
	}

	natSub, err = nats.NewPubSub(natsUrl, "mqtt", Logger)
	if err != nil {
		Logger.Error(fmt.Sprintf("Failed to connect to NATS: %s", err))
		return
	}
}
