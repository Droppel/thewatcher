package main

import (
	"watcher/archipelago"
	"watcher/discordbot"
)

func main() {
	err, messageCh := discordbot.InitBot()
	if err != nil {
		panic(err)
	}

	archipelago.Connect(messageCh)

}
