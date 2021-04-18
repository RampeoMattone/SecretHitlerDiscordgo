package bot

// This file describes the commands that will be subscribed by the bot

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "join", // join the game
			Description: "Join the lobby.",
		},
		{
			Name:        "start", // start the game
			Description: "Start the game.",
		},
		{
			Name:        "nominate", // nominate a chancellor
			Description: "Nominate a player as your chancellor candidate.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "nomination",
					Description: "The player you wish to nominate.",
					Required:    true,
				},
			},
		},
		{
			Name:        "vote", // vote for the goverment
			Description: "Vote to elect a government.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "Ja/Nein",
					Description: "Ja to say you're in favour of this election, Nein to say otherwise.",
					Required:    true,
				},
			},
		},
		{
			Name:        "remove_policy", // selects which policy to remove from the hand
			Description: "Select the policy you wish to remove from the group you've been dealt.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "selection",
					Description: "Select the policy to remove (0 is left, 1 is centre, 2 is right)",
					Required:    true,
				},
			},
		},
		{
			Name:        "veto", // ja or nein to the veto for a policy
			Description: "Vote to either veto the policy or pass it.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "Ja/Nein",
					Description: "Ja to say you agree to veto, Nein means otherwise.",
					Required:    true,
				},
			},
		},
		{
			Name:        "special", // special power
			Description: "Use your special power.",
			Options: []*discordgo.ApplicationCommandOption{
				//todo add subcommands
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
		// todo add command handlers
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
