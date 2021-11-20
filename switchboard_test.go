package switchboard

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
)

var testCommand = &discordgo.ApplicationCommand{
	Name:        "test",
	Description: "test",
}

var testCommandHandler = func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	return nil
}

func TestSwitchboardAddCommand_WithNewCommand(t *testing.T) {
	switchboard := &Switchboard{
		make(map[string]Command),
	}

	err := switchboard.AddCommand(testCommand, testCommandHandler)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if _, exists := switchboard.Commands[testCommand.Name]; !exists {
		t.Error("Command not added to Switchboard")
	}
}

func TestSwitchboardAddCommand_WithExistingCommand(t *testing.T) {
	switchboard := &Switchboard{
		make(map[string]Command),
	}

	switchboard.AddCommand(testCommand, testCommandHandler)
	err := switchboard.AddCommand(testCommand, testCommandHandler)

	if !errors.Is(err, ErrCommandAlreadyExists) {
		t.Errorf("Expected error %v, got %v", ErrCommandAlreadyExists, err)
	}
}

func TestSwitchboardHandleInteraction_WithUnsupportedInteractionType(t *testing.T) {
	switchboard := &Switchboard{
		make(map[string]Command),
	}

	err := switchboard.HandleInteractionCreate(nil, &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: 255,
		},
	})

	if !errors.Is(err, ErrUnsupportedInteractionType) {
		t.Errorf("Expected error %v, got %v", ErrUnsupportedInteractionType, err)
	}
}

func TestSwitchboardHandleInteraction_WithUnknownCommand(t *testing.T) {
	switchboard := &Switchboard{
		make(map[string]Command),
	}

	err := switchboard.HandleInteractionCreate(nil, &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test",
			},
		},
	})

	if !errors.Is(err, ErrUnknownCommand) {
		t.Errorf("Expected error %v, got %v", ErrUnknownCommand, err)
	}
}

func TestSwitchboardHandleInteraction_WithKnownCommand(t *testing.T) {
	switchboard := &Switchboard{
		make(map[string]Command),
	}

	called := false

	switchboard.AddCommand(testCommand, func(s *discordgo.Session, ic *discordgo.InteractionCreate) error {
		called = true
		return nil
	})

	err := switchboard.HandleInteractionCreate(nil, &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test",
			},
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !called {
		t.Error("Command handler not called")
	}
}

func TestSwitchboardHandlerInteraction_WithCommandErrorPropagates(t *testing.T) {
	switchboard := &Switchboard{
		make(map[string]Command),
	}

	expectedError := errors.New("test")

	switchboard.AddCommand(testCommand, func(s *discordgo.Session, ic *discordgo.InteractionCreate) error {
		return expectedError
	})

	err := switchboard.HandleInteractionCreate(nil, &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test",
			},
		},
	})

	if !errors.Is(err, expectedError) {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
}

func TestNew(t *testing.T) {
	switchboard := New()

	if switchboard.Commands == nil {
		t.Error("Expected Commands to be initialized, got nil")
	}
}
