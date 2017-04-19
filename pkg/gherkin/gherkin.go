package gherkin

import (
	"strings"

	irc "github.com/peted27/go-ircevent"
)

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
func MakeHelpString(command, info string) string {
	return command + ": " + info
}
