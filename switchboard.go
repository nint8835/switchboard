package switchboard

import "github.com/bwmarrin/discordgo"

// Switchboard is the main type responsible for managing all interactions.
type Switchboard struct {
}

func (switchboard *Switchboard) HandleInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	switch interaction.Type {
	default:
		return ErrUnsupportedInteractionType
	}
}

// New constructs a new Switchboard instance.
func New() *Switchboard {
	return &Switchboard{}
}
