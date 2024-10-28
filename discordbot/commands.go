package discordbot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
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

			currentGameStatus[channel.Name] = "Game status: BK"
			_, err = s.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
				Topic: "Game status: BK",
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

			currentGameStatus[channel.Name] = "Game status: SoftBK"
			s.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
				Topic: "Game status: SoftBK",
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

			currentGameStatus[channel.Name] = "Game status: unblocked"
			s.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
				Topic: "Game status: unblocked",
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
				case "Game status: BK":
					bkGames[chName] = topic
				case "Game status: SoftBK":
					softbkGames[chName] = topic
				case "Game status: unblocked":
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
