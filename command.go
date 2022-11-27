package switchboard

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Handler     any
	GuildID     string
}

func (c *Command) validate() error {
	err := validateHandler(c.Handler)

	if err != nil {
		return fmt.Errorf("invalid handler: %w", err)
	}

	return nil
}

func (c *Command) ToDiscordCommand() (*discordgo.ApplicationCommand, error) {
	err := c.validate()
	if err != nil {
		return nil, err
	}

	options, err := getCommandOptions(c.Handler)
	if err != nil {
		return nil, fmt.Errorf("error getting command options: %w", err)
	}

	return &discordgo.ApplicationCommand{
		Name:        c.Name,
		Description: c.Description,
		GuildID:     c.GuildID,
		Type:        discordgo.ChatApplicationCommand,
		Options:     options,
	}, nil
}
