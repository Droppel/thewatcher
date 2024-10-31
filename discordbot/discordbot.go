package discordbot

import (
	"fmt"
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

	currentGameStatus map[string]string = make(map[string]string)
)

type DiscordAction struct {
	Type             string
	Message          DiscordMessage
	ChannelTopicEdit DiscordChannelTopicEdit
}

type DiscordMessage struct {
	Slot    int
	Message string
	Silent  bool
}

type DiscordChannelTopicEdit struct {
	Slot  int
	Topic string
}

func InitBot() (chan DiscordAction, error) {
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

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		return nil, err
	}

	log.Info("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			log.Infof("Cannot create '%v' command: %v", v.Name, err)
			return nil, err
		}
		registeredCommands[i] = cmd
		fmt.Printf("Added command: %s\n", v.Name)
	}

	for _, chID := range slotsToChannels {
		channel, err := dg.Channel(chID)
		if err != nil {
			log.Println(err)
			continue
		}
		currentGameStatus[channel.Name] = channel.Topic
	}
	editStatusMessage()

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	messageCh := make(chan DiscordAction)

	go func() {
		for {
			select {
			case msg := <-messageCh:
				// Handle message
				switch msg.Type {
				case "message":
					var flags discordgo.MessageFlags = 0
					if msg.Message.Silent {
						flags = 4096
					}
					dg.ChannelMessageSendComplex(slotsToChannels[msg.Message.Slot], &discordgo.MessageSend{
						Content: msg.Message.Message,
						Flags:   flags,
					})
				case "channel_topic":
					channel, err := dg.Channel(slotsToChannels[msg.ChannelTopicEdit.Slot])
					if err != nil {
						log.Errorf("Cannot get channel: %v", err)
						continue
					}
					if !strings.Contains(channel.Topic, "BK") {
						continue
					}
					updateStatusMessage(channel.Name, msg.ChannelTopicEdit.Topic)
					dg.ChannelEdit(slotsToChannels[msg.ChannelTopicEdit.Slot], &discordgo.ChannelEdit{
						Topic: msg.ChannelTopicEdit.Topic,
					})
				}
			case <-sc:
				dg.Close()
				return
			}
		}
	}()

	return messageCh, err
}

func updateStatusMessage(gameName string, topic string) {
	currentGameStatus[gameName] = topic
	editStatusMessage()
}

func editStatusMessage() {
	channelID, exists := os.LookupEnv("STATUS_CHANNEL")
	if !exists {
		log.Error("STATUS_CHANNEL not set")
		return
	}

	channel, err := dg.Channel(channelID)
	if err != nil {
		log.Println(err)
		return
	}

	bkGames := make(map[string]string)
	softbkGames := make(map[string]string)
	unblockedGames := make(map[string]string)
	unknownGames := make(map[string]string)

	for chName, topic := range currentGameStatus {
		switch topic {
		case BK_STATUS:
			bkGames[chName] = topic
		case SOFTBK_STATUS:
			softbkGames[chName] = topic
		case UNBLOCKED_STATUS:
			unblockedGames[chName] = topic
		default:
			unknownGames[chName] = topic
		}
	}

	msgContent := "## Unknown games:\n"
	for chName, topic := range unknownGames {
		chReply := fmt.Sprintf("%s: %s\n", chName, topic)
		msgContent += chReply
	}

	msgContent += "\n## Unblocked games:\n"
	for chName, topic := range unblockedGames {
		chReply := fmt.Sprintf("%s: %s\n", chName, topic)
		msgContent += chReply
	}

	msgContent += "\n## SoftBK games:\n"
	for chName, topic := range softbkGames {
		chReply := fmt.Sprintf("%s: %s\n", chName, topic)
		msgContent += chReply
	}

	msgContent += "\n## BK games:\n"
	for chName, topic := range bkGames {
		chReply := fmt.Sprintf("%s: %s\n", chName, topic)
		msgContent += chReply
	}

	msgID := channel.LastMessageID
	_, err = dg.ChannelMessageEdit(channelID, msgID, msgContent)
	if err != nil {
		if strings.Contains(err.Error(), "Unknown Message") {
			dg.ChannelMessageSend(channelID, msgContent)
		}
		return
	}
}
