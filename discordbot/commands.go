package discordbot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"watcher/datastorage"

	"github.com/bwmarrin/discordgo"
)

const (
	BK_STATUS        = "Game status: BK"
	SOFTBK_STATUS    = "Game status: SoftBK"
	UNBLOCKED_STATUS = "Game status: unblocked"
	UNKNOWN_STATUS   = "Game status: unknown"
	GOAL_STATUS      = "Game status: goaled 🥳"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "bk",
			Description: "Sets the game status to BK",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "slot_number",
					Description: "Number of the slot",
					Required:    false,
				},
			},
		},
		{
			Name:        "softbk",
			Description: "Sets the game status to SoftBK",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "slot_number",
					Description: "Number of the slot",
					Required:    false,
				},
			},
		},
		{
			Name:        "unblocked",
			Description: "Sets the game status to unblocked",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "slot_number",
					Description: "Number of the slot",
					Required:    false,
				},
			},
		},
		{
			Name:        "goal",
			Description: "Sets the game status to goaled",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "slot_number",
					Description: "Number of the slot",
					Required:    false,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bk": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			updateStatusCommand(s, i, BK_STATUS)
		},
		"softbk": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			updateStatusCommand(s, i, SOFTBK_STATUS)
		},
		"unblocked": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			updateStatusCommand(s, i, UNBLOCKED_STATUS)
		},
		"goal": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			updateStatusCommand(s, i, GOAL_STATUS)
		},
	}
)

func updateStatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate, status string) {
	options := i.ApplicationCommandData().Options
	var slotIndex int64 = 1
	if len(options) > 0 {
		slotIndex = options[0].IntValue()
	}

	gameName := datastorage.SlotNumbersToAPSlots[channelsToSlots[i.ChannelID]].Name
	gameNameSplit := strings.Split(gameName, "_")
	gameName = gameNameSplit[0] // Remove the slot number
	maxSlot, _ := strconv.Atoi(gameNameSplit[1])
	if slotIndex < 1 {
		slotIndex = 1
	}
	if slotIndex > int64(maxSlot) {
		slotIndex = int64(maxSlot)
	}

	gameName = fmt.Sprintf("%s_%d", gameName, slotIndex)

	err := updateStatus(gameName, status)
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Failed to set game status to %s", status),
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Game status set to %s", status),
		},
	})
}

func init() {
}
