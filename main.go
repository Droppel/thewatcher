package main

import (
	"watcher/archipelago"
	"watcher/discordbot"
)

func main() {
	messageCh, err := discordbot.InitBot()
	if err != nil {
		panic(err)
	}

	archipelago.Connect(messageCh)
}
