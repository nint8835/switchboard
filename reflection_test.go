package switchboard

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-test/deep"
)

func Test_getCommandOptions_WithAllValidOptionTypes(t *testing.T) {
	minValue := 0.0

	options, err := getCommandOptions(
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct {
			String     string                      `description:"String argument"`
			Int        int                         `description:"Int argument"`
			Uint       uint                        `description:"Uint argument"`
			Bool       bool                        `description:"Bool argument"`
			User       discordgo.User              `description:"User argument"`
			Channel    discordgo.Channel           `description:"Channel argument"`
			Role       discordgo.Role              `description:"Role argument"`
			Float      float64                     `description:"Float argument"`
			Attachment discordgo.MessageAttachment `description:"Attachment argument"`
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
				Name:        "string",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionString,
				Description: "String argument",
			},
			{
				Name:        "int",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "Int argument",
			},
			{
				Name:        "uint",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "Uint argument",
				MinValue:    &minValue,
			},
			{
				Name:        "bool",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Description: "Bool argument",
			},
			{
				Name:        "user",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionUser,
				Description: "User argument",
			},
			{
				Name:        "channel",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionChannel,
				Description: "Channel argument",
			},
			{
				Name:        "role",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionRole,
				Description: "Role argument",
			},
			{
				Name:        "float",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionNumber,
				Description: "Float argument",
			},
			{
				Name:        "attachment",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Description: "Attachment argument",
			},
		},
	); diff != nil {
		t.Error(diff)
	}
}

func Test_getCommandOptions_WithPointerOption(t *testing.T) {
	options, err := getCommandOptions(
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct {
			Pointer *string `description:"Pointer"`
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
				Name:        "pointer",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
				Description: "Pointer",
			},
		},
	); diff != nil {
		t.Error(diff)
	}
}

func Test_getCommandOptions_WithDefaultValue(t *testing.T) {
	options, err := getCommandOptions(
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct {
			Default string `default:"default_val" description:"Default"`
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
				Name:        "default",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
				Description: "Default",
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

func Test_getCommandOptions_WithUintOption(t *testing.T) {
	minValue := 0.0

	options, err := getCommandOptions(
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct {
			Count uint `description:"Count argument"`
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
				Name:        "count",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "Count argument",
				MinValue:    &minValue,
			},
		},
	); diff != nil {
		t.Error(diff)
	}
}

func Test_getCommandOptions_WithUintPointerOption(t *testing.T) {
	minValue := 0.0

	options, err := getCommandOptions(
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args struct {
			Count *uint `description:"Count argument"`
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
				Name:        "count",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "Count argument",
				MinValue:    &minValue,
			},
		},
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

func Test_validateHandler_SlashCommand_WithNonFunction(t *testing.T) {
	err := validateHandler(SlashCommand, false)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerNotFunction) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_SlashCommand_WithWrongArgCount(t *testing.T) {
	err := validateHandler(SlashCommand, func() {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidParameterCount) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_SlashCommand_WithInvalidFirstArg(t *testing.T) {
	err := validateHandler(SlashCommand, func(first bool, second *discordgo.InteractionCreate, third struct{}) {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidFirstParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}

	err = validateHandler(
		SlashCommand,
		func(first discordgo.Session, second *discordgo.InteractionCreate, third struct{}) {}, //nolint:govet
	)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidFirstParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_SlashCommand_WithInvalidSecondArg(t *testing.T) {
	err := validateHandler(SlashCommand, func(first *discordgo.Session, second bool, third struct{}) {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidSecondParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}

	err = validateHandler(
		SlashCommand,
		func(first *discordgo.Session, second discordgo.InteractionCreate, third struct{}) {},
	)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidSecondParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_SlashCommand_WithInvalidThird(t *testing.T) {
	err := validateHandler(
		SlashCommand,
		func(first *discordgo.Session, second *discordgo.InteractionCreate, third bool) {},
	)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidThirdParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_SlashCommand_WithValidHandler(t *testing.T) {
	err := validateHandler(
		SlashCommand,
		func(first *discordgo.Session, second *discordgo.InteractionCreate, third struct{}) {},
	)

	if err != nil {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_MessageCommand_WithNonFunction(t *testing.T) {
	err := validateHandler(MessageCommand, false)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerNotFunction) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_MessageCommand_WithWrongArgCount(t *testing.T) {
	err := validateHandler(MessageCommand, func() {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidParameterCount) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_MessageCommand_WithInvalidFirstArg(t *testing.T) {
	err := validateHandler(
		MessageCommand,
		func(first bool, second *discordgo.InteractionCreate, third *discordgo.Message) {},
	)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidFirstParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}

	err = validateHandler(
		MessageCommand,
		func(first discordgo.Session, second *discordgo.InteractionCreate, third *discordgo.Message) {}, //nolint:govet
	)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidFirstParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_MessageCommand_WithInvalidSecondArg(t *testing.T) {
	err := validateHandler(MessageCommand, func(first *discordgo.Session, second bool, third *discordgo.Message) {})

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidSecondParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}

	err = validateHandler(
		MessageCommand,
		func(first *discordgo.Session, second discordgo.InteractionCreate, third *discordgo.Message) {},
	)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrHandlerInvalidSecondParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_MessageCommand_WithInvalidThird(t *testing.T) {
	err := validateHandler(
		MessageCommand,
		func(first *discordgo.Session, second *discordgo.InteractionCreate, third bool) {},
	)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrMessageHandlerInvalidThirdParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}

	err = validateHandler(
		MessageCommand,
		func(first *discordgo.Session, second *discordgo.InteractionCreate, third discordgo.Message) {},
	)

	if err == nil {
		t.Error("did not get expected error when validating handler")
	}
	if !errors.Is(err, ErrMessageHandlerInvalidThirdParameterType) {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_validateHandler_MessageCommand_WithValidHandler(t *testing.T) {
	err := validateHandler(
		MessageCommand,
		func(first *discordgo.Session, second *discordgo.InteractionCreate, third *discordgo.Message) {},
	)

	if err != nil {
		t.Errorf("got unexpected error when validating handler: %s", err)
	}
}

func Test_invokeCommand_SlashCommand_WithNoArgs(t *testing.T) {
	interactionData := discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{},
			},
		},
	}

	called := false

	invokeCommand(
		&Command{
			Type: SlashCommand,
		},
		&discordgo.Session{},
		&interactionData,
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, _ struct{}) {
			called = true
		})

	if !called {
		t.Error("handler function not called")
	}
}

func Test_invokeCommand_SlashCommand_WithProvidedArgs(t *testing.T) {
	interactionData := discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "string",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: "test",
					},
					{
						Name:  "int",
						Type:  discordgo.ApplicationCommandOptionInteger,
						Value: 12345.0,
					},
					{
						Name:  "bool",
						Type:  discordgo.ApplicationCommandOptionBoolean,
						Value: true,
					},
					{
						Name:  "user",
						Type:  discordgo.ApplicationCommandOptionUser,
						Value: "12345",
					},
					{
						Name:  "channel",
						Type:  discordgo.ApplicationCommandOptionChannel,
						Value: "12345",
					},
					{
						Name:  "role",
						Type:  discordgo.ApplicationCommandOptionRole,
						Value: "12345",
					},
					{
						Name:  "float",
						Type:  discordgo.ApplicationCommandOptionNumber,
						Value: 12345.0,
					},
					//{
					//	Name:  "attachment",
					//	Type:  discordgo.ApplicationCommandOptionAttachment,
					//	Value: "12345",
					//},
				},
			},
		},
	}

	called := false

	type Args struct {
		String  string
		Int     int
		Bool    bool
		User    discordgo.User
		Channel discordgo.Channel
		Role    discordgo.Role
		Float   float64
		//Attachment discordgo.MessageAttachment
	}

	invokeCommand(
		&Command{
			Type: SlashCommand,
		},
		nil,
		&interactionData,
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args Args) {
			called = true

			if diff := deep.Equal(
				args,
				Args{
					String:  "test",
					Int:     12345,
					Bool:    true,
					User:    discordgo.User{ID: "12345"},
					Channel: discordgo.Channel{ID: "12345"},
					Role:    discordgo.Role{ID: "12345"},
					Float:   12345.0,
					//Attachment: discordgo.MessageAttachment{ID: "12345"},
				},
			); diff != nil {
				t.Error(diff)
			}
		})

	if !called {
		t.Error("handler function not called")
	}
}

func Test_invokeCommand_SlashCommand_WithUintArg(t *testing.T) {
	interactionData := discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "count",
						Type:  discordgo.ApplicationCommandOptionInteger,
						Value: 42.0,
					},
				},
			},
		},
	}

	called := false

	type Args struct {
		Count uint
	}

	invokeCommand(
		&Command{
			Type: SlashCommand,
		},
		nil,
		&interactionData,
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, args Args) {
			called = true

			if diff := deep.Equal(
				args,
				Args{
					Count: 42,
				},
			); diff != nil {
				t.Error(diff)
			}
		})

	if !called {
		t.Error("handler function not called")
	}
}

func Test_invokeCommand_MessageCommand(t *testing.T) {
	interactionData := discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				TargetID: "1",

				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{Messages: map[string]*discordgo.Message{
					"1": {
						ID:      "1",
						Content: "Test",
						GuildID: "",
					},
				}},
			},
			GuildID: "2",
		},
	}

	called := false

	invokeCommand(
		&Command{
			Type: MessageCommand,
		},
		&discordgo.Session{},
		&interactionData,
		func(_ *discordgo.Session, _ *discordgo.InteractionCreate, msg *discordgo.Message) {
			called = true

			if diff := deep.Equal(
				msg,
				&discordgo.Message{
					ID:      "1",
					Content: "Test",
					GuildID: "2",
				},
			); diff != nil {
				t.Error(diff)
			}
		})

	if !called {
		t.Error("handler function not called")
	}
}
