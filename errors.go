package switchboard

import "errors"

// ErrCommandAlreadyExists occurs when a command is attmpting to be registered which has a name matching an already registered command.
var ErrCommandAlreadyExists = errors.New("command already exists")

// ErrUnknownCommand occurs when a command is attmpting to be handled which does not match any registered command..
var ErrUnknownCommand = errors.New("unknown command")

// ErrUnsupportedInteractionType occurs when an interaction is passed to *switchboard.HandleInteraction of a type that is not currently supported.
var ErrUnsupportedInteractionType = errors.New("unsupported interaction type")
