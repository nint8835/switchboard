package switchboard

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Switchboard struct {
	commands []*Command
}

func (s *Switchboard) handleInteractionApplicationCommand(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
) error {
	for _, command := range s.commands {
		if command.Name == interaction.ApplicationCommandData().Name &&
			(command.GuildID == "" || command.GuildID == interaction.GuildID) {
			invokeCommand(session, interaction, command.Handler)
			return nil
		}
	}

	return ErrUnknownCommand
}

func (s *Switchboard) HandleInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Type { //nolint:exhaustive
	case discordgo.InteractionApplicationCommand:
		s.handleInteractionApplicationCommand(session, interaction) //nolint:errcheck
	default:
		// TODO: Figure out error handling
	}
}

func (s *Switchboard) AddCommand(command *Command) error {
	// TODO: Add checks for things like duplicate commands, and ensure added command is valid
	s.commands = append(s.commands, command)

	return nil
}

func (s *Switchboard) RegisterCommands(session *discordgo.Session, appId string) error {
	for _, command := range s.commands {
		discordCmd, err := command.ToDiscordCommand()
		if err != nil {
			return fmt.Errorf("error generating discord command for command %s: %w", command.Name, err)
		}
		_, err = session.ApplicationCommandCreate(appId, command.GuildID, discordCmd)
		if err != nil {
			return fmt.Errorf("error registering discord command for command %s: %w", command.Name, err)
		}
	}

	return nil
}
