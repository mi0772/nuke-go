package tcp

import (
	"log"
	"strings"
)

type InputCommand struct {
	rawCommand        string
	commandIdentifier string
	key               string
	parameters        []string
}

func NewInputCommand(input string) InputCommand {
	parts := splitInput(input)

	cmd := InputCommand{
		rawCommand:        input,
		commandIdentifier: parts[0],
		key:               parts[1],
	}

	if len(parts) > 2 {
		cmd.parameters = parts[2:]
	}

	return cmd
}

func CommandBuilder(userCommand InputCommand) (Command, error) {
	return &PopCommand{}, nil
}

type Command interface {
	Process()
}

type PopCommand struct {
}

type PushCommand struct {
}

type ReadCommand struct {
}

func (c *PopCommand) Process() {
	log.Printf("ciao sono Process di GetCommand")
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
