package gherkin

import (
	"strings"

	"fmt"

	irc "github.com/peted27/go-ircevent"
)

type Plugin struct {
	Name    string
	Command string
	Help    string
	Version string
}

func hasCommandPrefix(s string) bool {
	if strings.HasPrefix(s, ":") || strings.HasPrefix(s, "!") {
		return true
	}
	return false
}

// IsPrivateMessage returns true if message e was private (ie not said in channel)
func IsPrivateMessage(e *irc.Event) bool {
	if strings.HasPrefix(e.Arguments[0], "#") {
		return false
	}
	return true

}

// IsCommandMessage returns true if message e was a command said in a channel, eg !ping !quit
func IsCommandMessage(e *irc.Event) bool {
	if strings.HasPrefix(e.Arguments[0], "#") && hasCommandPrefix(e.Arguments[1]) {
		return true
	}
	return false
}

// IsPublicMessage returns true if message e is a chat message said in channel
func IsPublicMessage(e *irc.Event) bool {
	if IsCommandMessage(e) || IsPrivateMessage(e) {
		return false
	}
	return true
}

// MakeHelpString creates a printable line for displaying help
func MakeHelpString(p Plugin) string {
	return fmt.Sprintf("%10s : %7s - %20s - %s", p.Name, p.Version, p.Command, p.Help)
}
