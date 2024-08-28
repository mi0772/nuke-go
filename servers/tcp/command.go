package tcp

import (
	"errors"
	"log"
	"strings"
)

type InputCommand struct {
	rawCommand        string
	commandIdentifier string
	key               string
	parameters        []string
}

func NewInputCommand(input string) (InputCommand, error) {
	parts := splitInput(input)

	// Controlla se ci sono abbastanza parti per formare un comando valido
	if len(parts) < 2 {
		return InputCommand{}, errors.New("invalid command")
	}

	cmd := InputCommand{
		rawCommand:        input,
		commandIdentifier: parts[0],
		key:               parts[1],
	}

	if len(parts) > 2 {
		cmd.parameters = parts[2:]
	}

	return cmd, nil
}

func CommandBuilder(userCommand InputCommand) (Command, error) {
	switch userCommand.commandIdentifier {
	case "POP":
		return &PopCommand{}, nil
	case "PUSH":
		return &PushCommand{}, nil
	case "READ":
		return &ReadCommand{}, nil
	default:
		return nil, errors.New("invalid command")
	}
}

type Command interface {
	Process(command *InputCommand)
}

type PopCommand struct {
}

type PushCommand struct {
}

type ReadCommand struct {
}

func (c *PopCommand) Process(command *InputCommand) {
	log.Printf("ciao sono Process di PopCommand : %s", command)
}

func (c *PushCommand) Process(command *InputCommand) {
	log.Printf("ciao sono Process di PushCommand : %s", command)
}

func (c *ReadCommand) Process(command *InputCommand) {
	log.Printf("ciao sono Process di ReadCommand : %s", command)
}

func splitInput(input string) []string {
	// This function splits the input string by spaces but preserves quoted substrings as single tokens
	var result []string
	var current strings.Builder
	var insideQuotes bool

	for _, char := range input {
		switch char {
		case ' ':
			if insideQuotes {
				current.WriteRune(char)
			} else if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		case '"':
			insideQuotes = !insideQuotes
		default:
			current.WriteRune(char)
		}
	}

	// Add the last part if there's anything left
	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
