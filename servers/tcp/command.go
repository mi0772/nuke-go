package tcp

import (
	"errors"
	"github.com/mi0772/nuke-go/engine"
	"github.com/mi0772/nuke-go/types"
	"strings"
	"time"
)

type CommandID int

const (
	Push   CommandID = 1
	Pop    CommandID = 2
	Read   CommandID = 3
	Quit   CommandID = 10
	Unknow CommandID = 999
)

func (id *CommandID) String() string {
	switch *id {
	case Push:
		return "PUSH"
	case Read:
		return "READ"
	case Pop:
		return "POP"
	case Quit:
		return "QUIT"
	default:
		return "UNKWNOW"
	}
}

func CommandIDFromString(s string) CommandID {
	switch strings.ToUpper(s) {
	case "PUSH":
		return Push
	case "POP":
		return Pop
	case "READ":
		return Read
	default:
		return Unknow

	}
}

type CommandInput struct {
	rawCommand        string
	commandIdentifier CommandID
	key               string
	value             []byte
}

func NewInputCommand(input string) (CommandInput, error) {
	if strings.Contains(input, "QUIT") {
		return CommandInput{
			commandIdentifier: CommandIDFromString("QUIT"),
		}, nil
	}
	parts := splitInput(input)

	// Controlla se ci sono abbastanza parti per formare un comando valido
	if len(parts) < 2 {
		return CommandInput{}, errors.New("invalid command")
	}

	cmd := CommandInput{
		rawCommand:        input,
		commandIdentifier: CommandIDFromString(parts[0]),
		key:               parts[1],
	}

	if len(parts) > 2 {
		if strings.HasPrefix(parts[2], "B:") {
			cmd.value = []byte(strings.TrimPrefix(parts[2], "B:"))
		} else {
			cmd.value = []byte(parts[2])
		}
	}

	return cmd, nil
}

func CommandBuilder(userCommand CommandInput) (Command, error) {
	switch userCommand.commandIdentifier {
	case Pop:
		return &PopCommand{}, nil
	case Push:
		return &PushCommand{}, nil
	case Read:
		return &ReadCommand{}, nil
	default:
		return nil, errors.New("invalid command")
	}
}

type Command interface {
	Process(command *CommandInput, database *engine.Database) (*engine.Item, types.NukeResponseCode)
}

type PopCommand struct {
}

type PushCommand struct {
}

type ReadCommand struct {
}

func (c *PopCommand) Process(command *CommandInput, database *engine.Database) (*engine.Item, types.NukeResponseCode) {
	start := time.Now()
	item, err := database.Pop(command.key)
	if err != nil {
		return nil, types.NOT_FOUND
	}
	end := time.Now()
	logf.Printf("Pop key:%s took %s\n", command.key, end.Sub(start))
	return &item, types.OK
}

func (c *PushCommand) Process(command *CommandInput, database *engine.Database) (*engine.Item, types.NukeResponseCode) {
	start := time.Now()
	item, err := database.Push(command.key, command.value)
	if err != nil {
		return nil, types.DUPLICATE_KEY
	}
	end := time.Now()
	logf.Printf("Push key:%s took %s", command.key, end.Sub(start))
	return &item, types.OK
}

func (c *ReadCommand) Process(command *CommandInput, database *engine.Database) (*engine.Item, types.NukeResponseCode) {
	start := time.Now()
	item, err := database.Read(command.key)
	if err != nil {
		return nil, types.NOT_FOUND
	}
	end := time.Now()
	logf.Printf("Read key:%s took %s", command.key, end.Sub(start))
	return item, types.OK
}

func splitInput(input string) []string {
	var result []string
	var current strings.Builder
	var insideQuotes bool

	for _, char := range input {
		switch char {
		case ';':
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

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
