package switchboard

import (
	"errors"
)

var ErrInvalidArgumentType = errors.New(
	"struct field type is not a supported option type",
)
var ErrHandlerNotFunction = errors.New(
	"provided command handler is not a function",
)
var ErrHandlerInvalidParameterCount = errors.New(
	"invalid number of handler arguments",
)
var ErrHandlerInvalidFirstParameterType = errors.New(
	"incorrect first parameter type for handler - first parameter must be of type *discordgo.Session",
)
var ErrHandlerInvalidSecondParameterType = errors.New(
	"incorrect second parameter type for handler - second parameter must be of type *discordgo.InteractionCreate",
)
var ErrHandlerInvalidThirdParameterType = errors.New(
	"incorrect third parameter type for handler - third parameter must be of type struct",
)
var ErrMessageHandlerInvalidThirdParameterType = errors.New(
	"incorrect third parameter type for handler - third parameter must be of type *discordgo.Message",
)
var ErrUnknownCommand = errors.New("unknown command")
var ErrUnsupportedInteractionType = errors.New("unsupported interaction type")
var ErrUnsupportedDefaultArgType = errors.New(
	"attempted to use default value for option type which does not currently support default values",
)
