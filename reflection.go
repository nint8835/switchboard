package switchboard

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var argTypeMap = map[reflect.Type]discordgo.ApplicationCommandOptionType{
	reflect.TypeOf(""):    discordgo.ApplicationCommandOptionString,
	reflect.TypeOf(0):     discordgo.ApplicationCommandOptionInteger,
	reflect.TypeOf(false): discordgo.ApplicationCommandOptionBoolean,
	// TODO: Should this be a `User` or a `Member`?
	reflect.TypeOf(discordgo.User{}):    discordgo.ApplicationCommandOptionUser,
	reflect.TypeOf(discordgo.Channel{}): discordgo.ApplicationCommandOptionChannel,
	reflect.TypeOf(discordgo.Role{}):    discordgo.ApplicationCommandOptionRole,
	// TODO: What type should `MENTIONABLE` map to?
	// ???: discordgo.ApplicationCommandOptionMentionable,
	reflect.TypeOf(0.0):                           discordgo.ApplicationCommandOptionNumber,
	reflect.TypeOf(discordgo.MessageAttachment{}): discordgo.ApplicationCommandOptionAttachment,
}

func getOptionType(argType reflect.Type) (discordgo.ApplicationCommandOptionType, error) {
	if argType.Kind() == reflect.Ptr {
		return getOptionType(argType.Elem())
	}

	argOptionType, validType := argTypeMap[argType]

	if validType {
		return argOptionType, nil
	} else {
		return 0, ErrInvalidArgumentType
	}
}

func getCommandOptions(argsStructType reflect.Type) ([]*discordgo.ApplicationCommandOption, error) {
	options := []*discordgo.ApplicationCommandOption{}

	for index := 0; index < argsStructType.NumField(); index++ {
		arg := argsStructType.Field(index)
		_, hasDefault := arg.Tag.Lookup("default")
		isPtr := arg.Type.Kind() == reflect.Ptr

		optionType, err := getOptionType(arg.Type)
		if err != nil {
			return nil, fmt.Errorf("unable to determine type for struct field %s: %w", arg.Name, err)
		}

		option := &discordgo.ApplicationCommandOption{
			Name:     strings.ToLower(arg.Name),
			Required: !(hasDefault || isPtr),
			Type:     optionType,
		}

		description, hasDescription := arg.Tag.Lookup("description")
		if hasDescription {
			option.Description = description
		}

		options = append(options, option)
	}

	return options, nil
}
