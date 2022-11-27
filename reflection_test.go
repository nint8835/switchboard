package switchboard

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-test/deep"
)

func Test_getCommandOptions_WithAllValidOptionTypes(t *testing.T) {
	options, err := getCommandOptions(
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct {
			String     string
			Int        int
			Bool       bool
			User       discordgo.User
			Channel    discordgo.Channel
			Role       discordgo.Role
			Float      float64
			Attachment discordgo.MessageAttachment
		}) {
		},
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
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct {
			Pointer *string
		}) {
		},
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
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct {
			Default string `default:"default_val"`
		}) {
		},
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
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct {
			Example string `description:"This is a test description"`
		}) {
		},
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
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct{}) {},
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
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct {
			Unsupported func()
		}) {
		},
	)

	if err == nil {
		t.Error("did not get expected error when getting command options")
	}
	if !errors.Is(err, ErrInvalidArgumentType) {
		t.Errorf("got unexpected error when getting command options: %s", err)
	}
}

func Test_validateHandler_WithNonFunction(t *testing.T) {
	err := validateHandler(false)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerNotFunction) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_WithWrongArgCount(t *testing.T) {
	err := validateHandler(func() {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidParameterCount) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_WithInvalidFirstArg(t *testing.T) {
	err := validateHandler(func(first bool, second *discordgo.InteractionCreate, third struct{}) {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidFirstParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}

	err = validateHandler(func(first discordgo.Session, second *discordgo.InteractionCreate, third struct{}) {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidFirstParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_WithInvalidSecondArg(t *testing.T) {
	err := validateHandler(func(first *discordgo.Session, second bool, third struct{}) {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidSecondParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}

	err = validateHandler(func(first *discordgo.Session, second discordgo.InteractionCreate, third struct{}) {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidSecondParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_WithInvalidThird(t *testing.T) {
	err := validateHandler(func(first *discordgo.Session, second *discordgo.InteractionCreate, third bool) {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidThirdParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_WithValidHandler(t *testing.T) {
	err := validateHandler(func(first *discordgo.Session, second *discordgo.InteractionCreate, third struct{}) {})

	if err != nil {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}
