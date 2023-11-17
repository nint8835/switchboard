package switchboard

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Switchboard struct {
	commands []*Command
}

func (s *Switchboard) hasCommand(name string, guildId string) bool {
	for _, command := range s.commands {
		if command.Name == name && command.GuildID == guildId {
			return true
		}
	}

	return false
}

func (s *Switchboard) handleInteractionApplicationCommand(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
) error {
	for _, command := range s.commands {
		if command.Name == interaction.ApplicationCommandData().Name &&
			(command.GuildID == "" || command.GuildID == interaction.GuildID) {
			invokeCommand(command, session, interaction, command.Handler)
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

func (s *Switchboard) SyncCommands(session *discordgo.Session, appId string) error {
	guildCommands := map[string][]*discordgo.ApplicationCommand{}

	for _, command := range s.commands {
		discordCommand, err := command.ToDiscordCommand()
		if err != nil {
			return fmt.Errorf("error generating discord command for command %s: %w", command.Name, err)
		}
		guildCommands[command.GuildID] = append(guildCommands[command.GuildID], discordCommand)
	}

	for guildId, commands := range guildCommands {
		_, err := session.ApplicationCommandBulkOverwrite(appId, guildId, commands)
		if err != nil {
			return fmt.Errorf("error syncing commands for guild %s: %w", guildId, err)
		}
	}

	return nil
}
