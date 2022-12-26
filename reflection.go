package switchboard

import (
	"fmt"
	"reflect"
	"strconv"
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

func getCommandOptions(handler any) ([]*discordgo.ApplicationCommandOption, error) {
	// Assumes validateHandler has been called before passing a handler to this function - will potentially panic otherwise
	argsStructType := reflect.TypeOf(handler).In(2)

	//goland:noinspection GoPreferNilSlice
	options := []*discordgo.ApplicationCommandOption{}

	for index := 0; index < argsStructType.NumField(); index++ {
		arg := argsStructType.Field(index)
		_, hasDefault := arg.Tag.Lookup("default")
		isPtr := arg.Type.Kind() == reflect.Ptr

		optionType, err := getOptionType(arg.Type)
		if err != nil {
			return nil, fmt.Errorf("unable to determine type for struct field %s: %w", arg.Name, err)
		}

		description, hasDescription := arg.Tag.Lookup("description")
		if !hasDescription {
			return nil, fmt.Errorf("no description provided for argument %s", arg.Name)
		}

		option := &discordgo.ApplicationCommandOption{
			Name:        strings.ToLower(arg.Name),
			Required:    !(hasDefault || isPtr),
			Type:        optionType,
			Description: description,
		}

		options = append(options, option)
	}

	return options, nil
}

func validateHandler(handler any) error {
	handlerType := reflect.TypeOf(handler)

	if handlerType.Kind() != reflect.Func {
		return ErrHandlerNotFunction
	}

	if handlerType.NumIn() != 3 {
		return ErrHandlerInvalidParameterCount
	}

	firstParam := handlerType.In(0)
	if firstParam.Kind() != reflect.Ptr || firstParam.Elem() != reflect.TypeOf(discordgo.Session{}) {
		return ErrHandlerInvalidFirstParameterType
	}

	secondParam := handlerType.In(1)
	if secondParam.Kind() != reflect.Ptr || secondParam.Elem() != reflect.TypeOf(discordgo.InteractionCreate{}) {
		return ErrHandlerInvalidSecondParameterType
	}

	if handlerType.In(2).Kind() != reflect.Struct {
		return ErrHandlerInvalidThirdParameterType
	}

	return nil
}

func getDefaultValue(field reflect.StructField) (reflect.Value, error) {
	defaultVal := field.Tag.Get("default")

	switch field.Type {
	case reflect.TypeOf(""):
		return reflect.ValueOf(defaultVal), nil
	case reflect.TypeOf(0):
		intVal, err := strconv.Atoi(defaultVal)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error parsing default value: %w", err)
		}
		return reflect.ValueOf(intVal), nil
	case reflect.TypeOf(false):
		boolVal, err := strconv.ParseBool(defaultVal)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error parsing default value: %w", err)
		}
		return reflect.ValueOf(boolVal), nil
	case reflect.TypeOf(0.0):
		floatVal, err := strconv.ParseFloat(defaultVal, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error parsing default value: %w", err)
		}
		return reflect.ValueOf(floatVal), nil
	default:
		return reflect.Value{}, ErrUnsupportedDefaultArgType
	}
}

func invokeCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate, handler any) {
	optionsMap := map[string]*discordgo.ApplicationCommandInteractionDataOption{}

	for _, option := range interaction.ApplicationCommandData().Options {
		optionsMap[option.Name] = option
	}

	argsParamType := reflect.TypeOf(handler).In(2)
	argsParamValue := reflect.New(argsParamType).Elem()

	for index := 0; index < argsParamValue.NumField(); index++ {
		field := argsParamValue.Field(index)
		fieldType := argsParamType.Field(index)

		option, optionProvided := optionsMap[strings.ToLower(fieldType.Name)]
		if optionProvided {
			var value reflect.Value

			switch option.Type { //nolint:exhaustive
			case discordgo.ApplicationCommandOptionString:
				value = reflect.ValueOf(option.StringValue())
			case discordgo.ApplicationCommandOptionInteger:
				value = reflect.ValueOf(int(option.IntValue()))
			case discordgo.ApplicationCommandOptionBoolean:
				value = reflect.ValueOf(option.BoolValue())
			// TODO: Is it fine to dereference users, roles, etc.?
			case discordgo.ApplicationCommandOptionUser:
				value = reflect.ValueOf(*option.UserValue(session))
			case discordgo.ApplicationCommandOptionChannel:
				value = reflect.ValueOf(*option.ChannelValue(session))
			case discordgo.ApplicationCommandOptionRole:
				value = reflect.ValueOf(*option.RoleValue(session, interaction.GuildID))
			case discordgo.ApplicationCommandOptionNumber:
				value = reflect.ValueOf(option.FloatValue())
				// TODO: find how to get attachment val
				//case discordgo.ApplicationCommandOptionAttachment:
				//	value = reflect.ValueOf(*option.)
			}

			if fieldType.Type.Kind() == reflect.Ptr {
				p := reflect.New(value.Type())
				p.Elem().Set(value)
				value = p
			}

			field.Set(value)
		} else if fieldType.Type.Kind() != reflect.Ptr {
			value, err := getDefaultValue(fieldType)
			if err != nil {
				// TODO: Handle properly, not panic
				panic(fmt.Errorf("error populating default value for field %s: %w", fieldType.Name, err))
			}

			field.Set(value)
		}
	}

	reflect.ValueOf(handler).Call(
		[]reflect.Value{
			reflect.ValueOf(session),
			reflect.ValueOf(interaction),
			argsParamValue,
		},
	)
}
