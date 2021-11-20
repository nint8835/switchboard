package switchboard

import "github.com/bwmarrin/discordgo"

type CommandHandler func(*discordgo.Session, *discordgo.InteractionCreate) error

// Command is a registered Switchboard command.
type Command struct {
	*discordgo.ApplicationCommand
	Handler CommandHandler
}

// Switchboard is the main type responsible for managing all interactions.
type Switchboard struct {
	Commands map[string]Command
}

// Register a command with Switchboard. Note that this does not register the command with Discord.
func (switchboard *Switchboard) AddCommand(command *discordgo.ApplicationCommand, handler CommandHandler) error {
	if _, exists := switchboard.Commands[command.Name]; exists {
		return ErrCommandAlreadyExists
	}

	switchboard.Commands[command.Name] = Command{command, handler}

	return nil
}

func (switchboard *Switchboard) handleInteractionApplicationCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	command, exists := switchboard.Commands[interaction.ApplicationCommandData().Name]
	if !exists {
		return ErrUnknownCommand
	}

	return command.Handler(session, interaction)
}

func (switchboard *Switchboard) HandleInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		return switchboard.handleInteractionApplicationCommand(session, interaction)
	default:
		return ErrUnsupportedInteractionType
	}
}

// New constructs a new Switchboard instance.
func New() *Switchboard {
	return &Switchboard{
		Commands: make(map[string]Command),
	}
}
