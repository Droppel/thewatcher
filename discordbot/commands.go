package discordbot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	BK_STATUS        = "Game status: BK"
	SOFTBK_STATUS    = "Game status: SoftBK"
	UNBLOCKED_STATUS = "Game status: unblocked"
	UNKNOWN_STATUS   = "Game status: unknown"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "bk",
			Description: "Sets the game status to BK",
		},
		{
			Name:        "softbk",
			Description: "Sets the game status to SoftBK",
		},
		{
			Name:        "unblocked",
			Description: "Sets the game status to unblocked",
		},
		{
			Name:        "bkstatus",
			Description: "Replies with the current game status",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bk": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				log.Println(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to get channel",
					},
				})
				return
			}

			updateStatusMessage(channel.Name, BK_STATUS)
			_, err = s.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
				Topic: BK_STATUS,
			})
			if err != nil {
				log.Println(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to set game status to BK",
					},
				})
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Game status set to BK",
				},
			})
		},
		"softbk": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				log.Println(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to get channel",
					},
				})
				return
			}

			updateStatusMessage(channel.Name, SOFTBK_STATUS)
			s.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
				Topic: SOFTBK_STATUS,
			})

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Game status set to SoftBK",
				},
			})
		},
		"unblocked": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				log.Println(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to get channel",
					},
				})
				return
			}

			updateStatusMessage(channel.Name, UNBLOCKED_STATUS)
			s.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
				Topic: UNBLOCKED_STATUS,
			})

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Game status set to unblocked",
				},
			})
		},
		"bkstatus": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

			reply := "## Unknown games:\n"
			for chName, topic := range unknownGames {
				chReply := fmt.Sprintf("%s: %s\n", chName, topic)
				reply += chReply
			}

			reply += "\n## Unblocked games:\n"
			for chName, topic := range unblockedGames {
				chReply := fmt.Sprintf("%s: %s\n", chName, topic)
				reply += chReply
			}

			reply += "\n## SoftBK games:\n"
			for chName, topic := range softbkGames {
				chReply := fmt.Sprintf("%s: %s\n", chName, topic)
				reply += chReply
			}

			reply += "\n## BK games:\n"
			for chName, topic := range bkGames {
				chReply := fmt.Sprintf("%s: %s\n", chName, topic)
				reply += chReply
			}

			log.Println("Replying with:", reply)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: reply,
				},
			})
		},
	}
)

func init() {
}
