package lib

import (
	"strings"

	"github.com/peted27/go-ircevent"
)

func hasDirectPrefix(s string) bool {
	if strings.HasPrefix(s, ":") || strings.HasPrefix(s, "!") {
		return true
	}
	return false
}

func IsPrivateMessage(e *irc.Event) bool {
	if strings.HasPrefix(e.Arguments[0], "#") {
		return false
	}
	return true

}

func IsDirectMessage(e *irc.Event) bool {
	if strings.HasPrefix(e.Arguments[0], "#") && hasDirectPrefix(e.Arguments[1]) {
		return true
	}
	return false
}

func IsCommandMessage(e *irc.Event) bool {
	if strings.HasPrefix(e.Arguments[0], "#") && strings.HasPrefix(e.Arguments[1], "!") {
		return true
	}
	return false
}

func IsChatMessage(e *irc.Event) bool {
	if IsCommandMessage(e) || IsDirectMessage(e) || IsPrivateMessage(e) {
		return false
	}
	return true
}
