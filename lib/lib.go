package lib

import (
	"strings"

	"github.com/peted27/go-ircevent"
)

func IsPrivate(e *irc.Event) bool {
	if strings.HasPrefix(e.Arguments[0], "#") {
		return false
	}
	return true

}

func hasDirectPrefix(s string) bool {
	if strings.HasPrefix(s, ":") || strings.HasPrefix(s, "!") {
		return true
	}
	return false
}

func IsDirect(e *irc.Event) bool {
	if strings.HasPrefix(e.Arguments[0], "#") && hasDirectPrefix(e.Arguments[1]) {
		return true
	}
	return false
}

func IsCommand(e *irc.Event) bool {
	if strings.HasPrefix(e.Arguments[0], "#") && strings.HasPrefix(e.Arguments[1], "!") {
		return true
	}
	return false
}
