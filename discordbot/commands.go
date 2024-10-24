package discordbot

import (
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
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bk": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			_, err := s.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
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
	}
)

func init() {
}
