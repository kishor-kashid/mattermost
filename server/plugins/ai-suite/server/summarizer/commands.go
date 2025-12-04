package summarizer

import (
	"fmt"
	"strings"
)

// CommandTrigger is the slash command keyword (without the slash).
const CommandTrigger = "summarize"

// CommandOptions contains parsed slash command arguments.
type CommandOptions struct {
	Target   Type
	Argument string
}

// ParseCommand parses `/summarize ...` invocations.
func ParseCommand(input string) (CommandOptions, error) {
	fields := strings.Fields(strings.TrimSpace(input))
	if len(fields) == 0 {
		return CommandOptions{}, fmt.Errorf("%w: missing command", ErrInvalidRequest)
	}

	trigger := strings.TrimPrefix(fields[0], "/")
	if strings.ToLower(trigger) != CommandTrigger {
		return CommandOptions{}, fmt.Errorf("%w: unsupported command %q", ErrInvalidRequest, trigger)
	}

	if len(fields) == 1 {
		return CommandOptions{}, fmt.Errorf("%w: specify 'thread' or 'channel'", ErrInvalidRequest)
	}

	target := Type(strings.ToLower(fields[1]))
	if target != TypeThread && target != TypeChannel {
		return CommandOptions{}, fmt.Errorf("%w: expected 'thread' or 'channel'", ErrInvalidRequest)
	}

	arg := ""
	if len(fields) > 2 {
		arg = strings.Join(fields[2:], " ")
	}

	return CommandOptions{
		Target:   target,
		Argument: strings.TrimSpace(arg),
	}, nil
}
