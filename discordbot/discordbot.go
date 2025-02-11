package discordbot

import (
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"watcher/datastorage"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	dg              *discordgo.Session
	slotsToChannels map[int]string
	// Because we have a 1 to N relation between channels and slots, this map is not complete, it doesn't matter though, because we only care
	// about any slot name and strip the number anyways
	channelsToSlots map[string]int

	statusMessageChannelID string
	statusMessageID        string

	CurrentGameStatus map[string]string = make(map[string]string)
)

type DiscordAction struct {
	Type         string
	Message      DiscordMessage
	StatusChange DiscordStatusChange
}

type DiscordMessage struct {
	Slot    int
	Message string
	Silent  bool
}

type DiscordStatusChange struct {
	Name   string
	Slot   int
	Status string
}

func InitBot() (chan DiscordAction, error) {
	var err error
	authtoken := os.Getenv("AUTH_TOKEN")

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

	statusMessageChannelID, statusMessageID, err = getStatusMessage()
	if err != nil {
		log.Error(err)
		return nil, err
	}

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

					if CurrentGameStatus[datastorage.SlotNumbersToAPSlots[msg.Message.Slot].Name] == GOAL_STATUS {
						continue
					}

					dg.ChannelMessageSendComplex(slotsToChannels[msg.Message.Slot], &discordgo.MessageSend{
						Content: msg.Message.Message,
						Flags:   flags,
					})
				case "status_change":
					skip_status := []string{GOAL_STATUS, UNKNOWN_STATUS, UNBLOCKED_STATUS}

					if slices.Contains[[]string, string](skip_status, CurrentGameStatus[datastorage.SlotNumbersToAPSlots[msg.Message.Slot].Name]) {
						continue
					}

					err = updateStatus(msg.StatusChange.Name, msg.StatusChange.Status)
					if err != nil {
						log.Errorf("Cannot edit channel: %v", err)
					}
				}
			case <-sc:
				dg.Close()
				return
			}
		}
	}()

	return messageCh, err
}

func InitBotAfterAPConnect() error {

	slotsToChannelsEnv := os.Getenv("SLOTS_TO_CHANNELS")
	slotsToChannels = make(map[int]string)
	channelsToSlots = make(map[string]int)

	for _, slotToChannel := range strings.Split(slotsToChannelsEnv, ",") {
		slotToChannelSplit := strings.Split(slotToChannel, ":")
		number, _ := strconv.Atoi(slotToChannelSplit[0])
		slotsToChannels[number] = slotToChannelSplit[1]
		channelsToSlots[slotToChannelSplit[1]] = number
	}

	statusMsg, err := dg.ChannelMessage(statusMessageChannelID, statusMessageID)
	if err == nil {
		currentStatus := BK_STATUS
		for _, msg := range strings.Split(statusMsg.Content, "\n") {
			if len(msg) == 0 {
				continue
			}
			if strings.Contains(msg, "##") {
				switch msg {
				case "## Unknown games:":
					currentStatus = UNKNOWN_STATUS
				case "## Unblocked games:":
					currentStatus = UNBLOCKED_STATUS
				case "## SoftBK games:":
					currentStatus = SOFTBK_STATUS
				case "## BK games:":
					currentStatus = BK_STATUS
				case "## Goaled games:":
					currentStatus = GOAL_STATUS
				}
				continue
			}
			CurrentGameStatus[msg] = currentStatus
		}
	}
	err = editStatusMessage()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func updateStatus(name string, newStatus string) error {
	if status, exists := CurrentGameStatus[name]; exists {
		if status == newStatus {
			return nil
		}
	}
	CurrentGameStatus[name] = newStatus
	return editStatusMessage()
}

func getStatusMessage() (string, string, error) {
	channelID, exists := os.LookupEnv("STATUS_CHANNEL")
	if !exists {
		return "", "", fmt.Errorf("STATUS_CHANNEL not set")
	}

	channel, err := dg.Channel(channelID)
	if err != nil {
		return "", "", err
	}

	return channel.ID, channel.LastMessageID, nil
}

func editStatusMessage() error {
	channelID, statusMsgID, err := getStatusMessage()
	if err != nil {
		return err
	}

	bkGames := make([]string, 0)
	softbkGames := make([]string, 0)
	unblockedGames := make([]string, 0)
	unknownGames := make([]string, 0)
	goaledGames := make([]string, 0)

	for name, status := range CurrentGameStatus {
		switch status {
		case BK_STATUS:
			bkGames = append(bkGames, name)
		case SOFTBK_STATUS:
			softbkGames = append(softbkGames, name)
		case UNBLOCKED_STATUS:
			unblockedGames = append(unblockedGames, name)
		case GOAL_STATUS:
			goaledGames = append(goaledGames, name)
		default:
			unknownGames = append(unknownGames, name)
		}
	}

	slices.Sort(bkGames)
	slices.Sort(softbkGames)
	slices.Sort(unblockedGames)
	slices.Sort(unknownGames)
	slices.Sort(goaledGames)

	msgContent := "## Unknown games:\n"
	for _, name := range unknownGames {
		chReply := fmt.Sprintf("%s\n", name)
		msgContent += chReply
	}

	msgContent += "\n## Unblocked games:\n"
	for _, name := range unblockedGames {
		chReply := fmt.Sprintf("%s\n", name)
		msgContent += chReply
	}

	msgContent += "\n## SoftBK games:\n"
	for _, name := range softbkGames {
		chReply := fmt.Sprintf("%s\n", name)
		msgContent += chReply
	}

	msgContent += "\n## BK games:\n"
	for _, name := range bkGames {
		chReply := fmt.Sprintf("%s\n", name)
		msgContent += chReply
	}

	msgContent += "\n## Goaled games:\n"
	for _, name := range goaledGames {
		chReply := fmt.Sprintf("%s\n", name)
		msgContent += chReply
	}

	_, err = dg.ChannelMessageEdit(channelID, statusMsgID, msgContent)
	if err != nil {
		if strings.Contains(err.Error(), "Unknown Message") {
			_, err = dg.ChannelMessageSend(channelID, msgContent)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}
