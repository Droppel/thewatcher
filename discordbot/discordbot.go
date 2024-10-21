package discordbot

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	dg              *discordgo.Session
	slotsToChannels map[int]string
)

type DiscordMessage struct {
	Slot    int
	Message string
}

func InitBot() (chan DiscordMessage, error) {
	var err error
	authtoken := os.Getenv("AUTH_TOKEN")

	slotsToChannelsEnv := os.Getenv("SLOTS_TO_CHANNELS")
	slotsToChannels = make(map[int]string)
	for _, slotToChannel := range strings.Split(slotsToChannelsEnv, ",") {
		slotToChannelSplit := strings.Split(slotToChannel, ":")
		number, _ := strconv.Atoi(slotToChannelSplit[0])
		slotsToChannels[number] = slotToChannelSplit[1]
	}

	// Create a new Discord session using the provided bot token.
	dg, err = discordgo.New("Bot " + string(authtoken))
	if err != nil {
		return nil, err
	}

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		return nil, err
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	messageCh := make(chan DiscordMessage)

	go func() {
		for {
			select {
			case msg := <-messageCh:
				// Handle message

				dg.ChannelMessageSend(slotsToChannels[msg.Slot], msg.Message)
			case <-sc:
				dg.Close()
				return
			}
		}
	}()

	return messageCh, err
}
