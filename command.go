package switchboard

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Handler     any
	GuildID     string
}

func (c *Command) validate() error {
	return nil
}

func (c *Command) toDiscordCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name,
		Description: c.Description,
		GuildID:     c.GuildID,
		Type:        discordgo.ChatApplicationCommand,
	}
}
