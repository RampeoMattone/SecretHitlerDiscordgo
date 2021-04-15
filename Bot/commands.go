package bot

// This file describes the commands that will be subscribed by the bot

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "admin",
			Description: "pong a ping",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "delay",
					Description: "add a delay to the response",
					Required:    false,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"are_you_there": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			if err != nil {
				return
			}
			for _, option := range i.Data.Options {
				if option.Name == "delay" && option.Value == true {
					time.Sleep(10 * time.Second)
				}
			}
			err = s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
				Content: "Hey there! I can confirm that I'm alive <3",
			})
			if err != nil {
				return
			}
		},
	}
)
