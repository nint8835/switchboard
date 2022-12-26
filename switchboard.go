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

func (s *Switchboard) SyncCommands(session *discordgo.Session, appId string) error {
	// TODO: There's probably a cleaner way to do all this - check how discord.py handles it
	globalCommands, err := session.ApplicationCommands(appId, "")
	if err != nil {
		return fmt.Errorf("error listing global commands: %w", err)
	}

	for _, globalCommand := range globalCommands {
		if !s.hasCommand(globalCommand.Name, globalCommand.GuildID) {
			err = session.ApplicationCommandDelete(appId, "", globalCommand.ID)
			if err != nil {
				return fmt.Errorf("error deleting command %s: %w", globalCommand.Name, err)
			}
		}
	}

	guilds, err := session.UserGuilds(100, "", "")
	if err != nil {
		return fmt.Errorf("error listing guilds: %w", err)
	}

	for _, guild := range guilds {
		guildCommands, err := session.ApplicationCommands(appId, guild.ID)
		if err != nil {
			return fmt.Errorf("error listing commands for guild %s: %w", guild.ID, err)
		}

		for _, guildCommand := range guildCommands {
			if !s.hasCommand(guildCommand.Name, guildCommand.GuildID) {
				err = session.ApplicationCommandDelete(appId, guildCommand.GuildID, guildCommand.ID)
				if err != nil {
					return fmt.Errorf("error deleting command %s: %w", guildCommand.Name, err)
				}
			}
		}
	}

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
