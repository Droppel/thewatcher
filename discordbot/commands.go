package discordbot

import (
	"log"

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
			Name:        "goal",
			Description: "Sets the game status to goaled",
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
		"goal": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

			updateStatusMessage(channel.Name, GOAL_STATUS)
			s.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
				Topic: GOAL_STATUS,
			})

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Game status set to goaled",
				},
			})
		},
	}
)

func init() {
}
