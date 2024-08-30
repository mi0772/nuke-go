package tcp

import (
	"encoding/json"
	"errors"
	"github.com/mi0772/nuke-go/engine"
	"github.com/mi0772/nuke-go/types"
	"strings"
	"time"
)

type CommandID int

const (
	Push             CommandID = 1
	Pop              CommandID = 2
	Read             CommandID = 3
	Keys             CommandID = 4
	PartitionDetails CommandID = 5
	Quit             CommandID = 10
	Unknow           CommandID = 999
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
	case Keys:
		return "KEYS"
	case PartitionDetails:
		return "PARTITION_DETAILS"
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
	case "QUIT":
		return Quit
	case "KEYS":
		return Keys
	case "PD":
		return PartitionDetails
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
	// Maps with commands that don't need a key
	commands := map[string]string{
		"QUIT": "QUIT",
		"KEYS": "KEYS",
		"PD":   "PD",
	}

	for key, commandID := range commands {
		if strings.Contains(input, key) {
			return CommandInput{
				commandIdentifier: CommandIDFromString(commandID),
			}, nil
		}
	}

	//if we are here, we have a command that needs a key
	parts := splitInput(input)

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
	case Keys:
		return &KeysListCommand{}, nil
	case PartitionDetails:
		return &PartitionDetailsCommand{}, nil

	default:
		return nil, errors.New("invalid command")
	}
}

type Command interface {
	Process(command *CommandInput, database *engine.Database) ([]byte, types.NukeResponseCode)
}

type PopCommand struct {
}

type PushCommand struct {
}

type ReadCommand struct {
}

type KeysListCommand struct {
}

type PartitionDetailsCommand struct {
}

func (c *KeysListCommand) Process(command *CommandInput, database *engine.Database) ([]byte, types.NukeResponseCode) {
	start := time.Now()
	keys := database.Keys()
	end := time.Now()
	logf.Printf("KeysList took %s\n", end.Sub(start))
	return toJSON(keys), types.Ok
}

func (c *PartitionDetailsCommand) Process(command *CommandInput, database *engine.Database) ([]byte, types.NukeResponseCode) {
	start := time.Now()
	partitions := database.DetailsPartitions()
	end := time.Now()
	logf.Printf("PartitionDetails took %s\n", end.Sub(start))
	return toJSON(partitions), types.Ok
}

func (c *PopCommand) Process(command *CommandInput, database *engine.Database) ([]byte, types.NukeResponseCode) {
	start := time.Now()
	item, err := database.Pop(command.key)
	if err != nil {
		return nil, types.NotFound
	}
	end := time.Now()
	logf.Printf("Pop key:%s took %s\n", command.key, end.Sub(start))
	return toJSON(item), types.Ok
}

func (c *PushCommand) Process(command *CommandInput, database *engine.Database) ([]byte, types.NukeResponseCode) {
	start := time.Now()
	item, err := database.Push(command.key, command.value)
	if err != nil {
		return nil, types.DuplicateKey
	}
	end := time.Now()
	logf.Printf("Push key:%s took %s", command.key, end.Sub(start))
	return toJSON(item), types.Ok
}

func (c *ReadCommand) Process(command *CommandInput, database *engine.Database) ([]byte, types.NukeResponseCode) {
	start := time.Now()
	item, err := database.Read(command.key)
	if err != nil {
		return nil, types.NotFound
	}
	end := time.Now()
	logf.Printf("Read key:%s took %s", command.key, end.Sub(start))
	return toJSON(item), types.Ok
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

func toJSON(c interface{}) []byte {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return []byte{}
	}
	return jsonData
}
