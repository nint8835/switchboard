package switchboard

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type CommandType int

const (
	SlashCommand CommandType = iota
	MessageCommand
)

var typeMap = map[CommandType]discordgo.ApplicationCommandType{
	SlashCommand:   discordgo.ChatApplicationCommand,
	MessageCommand: discordgo.MessageApplicationCommand,
}

type Command struct {
	Name        string
	Description string
	Handler     any
	GuildID     string

	Type CommandType
}

func (c *Command) validate() error {
	err := validationFuncs[c.Type](c.Handler)

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

	var options []*discordgo.ApplicationCommandOption

	if c.Type == SlashCommand {
		options, err = getCommandOptions(c.Handler)
		if err != nil {
			return nil, fmt.Errorf("error getting command options: %w", err)
		}
	}

	return &discordgo.ApplicationCommand{
		Name:        c.Name,
		Description: c.Description,
		GuildID:     c.GuildID,
		Type:        typeMap[c.Type],
		Options:     options,
	}, nil
}
