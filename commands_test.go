package switchboard

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-test/deep"
)

func Test_getCommandOptions_WithAllValidOptionTypes(t *testing.T) {
	options, err := getCommandOptions(
		reflect.TypeOf(
			struct {
				String     string
				Int        int
				Bool       bool
				User       discordgo.User
				Channel    discordgo.Channel
				Role       discordgo.Role
				Float      float64
				Attachment discordgo.MessageAttachment
			}{},
		),
	)
	if err != nil {
		t.Errorf("got unexpected error getting command options: %s", err)
	}

	if diff := deep.Equal(
		options,
		[]*discordgo.ApplicationCommandOption{
			{
				Name:     "string",
				Required: true,
				Type:     discordgo.ApplicationCommandOptionString,
			},
			{
				Name:     "int",
				Required: true,
				Type:     discordgo.ApplicationCommandOptionInteger,
			},
			{
				Name:     "bool",
				Required: true,
				Type:     discordgo.ApplicationCommandOptionBoolean,
			},
			{
				Name:     "user",
				Required: true,
				Type:     discordgo.ApplicationCommandOptionUser,
			},
			{
				Name:     "channel",
				Required: true,
				Type:     discordgo.ApplicationCommandOptionChannel,
			},
			{
				Name:     "role",
				Required: true,
				Type:     discordgo.ApplicationCommandOptionRole,
			},
			{
				Name:     "float",
				Required: true,
				Type:     discordgo.ApplicationCommandOptionNumber,
			},
			{
				Name:     "attachment",
				Required: true,
				Type:     discordgo.ApplicationCommandOptionAttachment,
			},
		},
	); diff != nil {
		t.Error(diff)
	}
}

func Test_getCommandOptions_WithPointerOption(t *testing.T) {
	options, err := getCommandOptions(
		reflect.TypeOf(struct {
			Pointer *string
		}{}),
	)

	if err != nil {
		t.Errorf("got unexpected error getting command options: %s", err)
	}

	if diff := deep.Equal(
		options,
		[]*discordgo.ApplicationCommandOption{
			{
				Name:     "pointer",
				Required: false,
				Type:     discordgo.ApplicationCommandOptionString,
			},
		},
	); diff != nil {
		t.Error(diff)
	}
}

func Test_getCommandOptions_WithDefaultValue(t *testing.T) {
	options, err := getCommandOptions(
		reflect.TypeOf(struct {
			Default string `default:"default_val"`
		}{}),
	)

	if err != nil {
		t.Errorf("got unexpected error getting command options: %s", err)
	}

	if diff := deep.Equal(
		options,
		[]*discordgo.ApplicationCommandOption{
			{
				Name:     "default",
				Required: false,
				Type:     discordgo.ApplicationCommandOptionString,
			},
		},
	); diff != nil {
		t.Error(diff)
	}
}

func Test_getCommandOptions_WithDescription(t *testing.T) {
	options, err := getCommandOptions(
		reflect.TypeOf(struct {
			Example string `description:"This is a test description"`
		}{}),
	)

	if err != nil {
		t.Errorf("got unexpected error getting command options: %s", err)
	}

	if diff := deep.Equal(
		options,
		[]*discordgo.ApplicationCommandOption{
			{
				Name:        "example",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionString,
				Description: "This is a test description",
			},
		},
	); diff != nil {
		t.Error(diff)
	}
}

func Test_getCommandOptions_WithNoOptions(t *testing.T) {
	options, err := getCommandOptions(
		reflect.TypeOf(struct{}{}),
	)

	if err != nil {
		t.Errorf("got unexpected error getting command options: %s", err)
	}

	if diff := deep.Equal(
		options,
		[]*discordgo.ApplicationCommandOption{},
	); diff != nil {
		t.Error(diff)
	}
}

func Test_getCommandOptions_WithInvalidArgumentType(t *testing.T) {
	_, err := getCommandOptions(
		reflect.TypeOf(struct {
			Unsupported func()
		}{}),
	)

	if err == nil {
		t.Error("did not get expected error when getting command options")
	}
	if !errors.Is(err, ErrInvalidArgumentType) {
		t.Errorf("got unexpected error when getting command options: %s", err)
	}
}
