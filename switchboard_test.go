package switchboard

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestSwitchboardHandleInteraction_WithUnsupportedInteractionType(t *testing.T) {
	switchboard := &Switchboard{}

	err := switchboard.HandleInteractionCreate(nil, &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: 255,
		},
	})

	if !errors.Is(err, ErrUnsupportedInteractionType) {
		t.Errorf("Expected error %v, got %v", ErrUnsupportedInteractionType, err)
	}
}
