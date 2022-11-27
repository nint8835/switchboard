package switchboard

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-test/deep"
)

func TestCommand_ToDiscordCommand_WithHandlerError(t *testing.T) {
	testCommand := &Command{
		Name:        "",
		Description: "",
		Handler:     false,
		GuildID:     "",
	}

	_, err := testCommand.ToDiscordCommand()

	if err == nil {
		t.Error("did not get expected error")
	}
	if !errors.Is(err, ErrHandlerNotFunction) {
		t.Errorf("got unexpected error: %s", err)
	}
}

func TestCommand_ToDiscordCommand_WithOptionError(t *testing.T) {
	testCommand := &Command{
		Name:        "",
		Description: "",
		Handler: func(
			_ *discordgo.Session,
			_ *discordgo.InteractionCreate,
			args struct {
				Arg map[string]any
			},
		) {

		},
		GuildID: "",
	}

	_, err := testCommand.ToDiscordCommand()

	if err == nil {
		t.Error("did not get expected error")
	}
	if !errors.Is(err, ErrInvalidArgumentType) {
		t.Errorf("got unexpected error: %s", err)
	}
}

func TestCommand_ToDiscordCommand_WithValidCommand(t *testing.T) {
	testCommand := &Command{
		Name:        "test",
		Description: "This is a test command",
		Handler: func(
			_ *discordgo.Session,
			_ *discordgo.InteractionCreate,
			args struct {
				Required string  `description:"This is a required arg"`
				Opt1     int     `description:"This is an optional arg" default:"5"`
				Opt2     *string `description:"This is an optional arg as well"`
			},
		) {

		},
		GuildID: "1234567890",
	}

	cmd, err := testCommand.ToDiscordCommand()

	if err != nil {
		t.Errorf("got unexpected error: %s", err)
	}

	if diff := deep.Equal(
		cmd,
		&discordgo.ApplicationCommand{
			Name:        "test",
			Description: "This is a test command",
			Type:        discordgo.ChatApplicationCommand,
			GuildID:     "1234567890",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "required",
					Description: "This is a required arg",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "opt1",
					Description: "This is an optional arg",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
				},
				{
					Name:        "opt2",
					Description: "This is an optional arg as well",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
			},
		},
	); diff != nil {
		t.Error(diff)
	}
}
