package lib

import (
	"strings"

	"github.com/peted27/go-ircevent"
)

func PrivmsgHandler(f func(e *irc.Event), hDirect, hIndirect, hPrivate bool) func(e *irc.Event) {

	return func(e *irc.Event) {
		var private, direct, indirect bool

		// deterimine the type of private message
		if strings.HasPrefix(e.Arguments[0], "#") {
			switch {
			case strings.HasPrefix(e.Arguments[1], "!"):
				indirect = true
			case strings.HasPrefix(e.Arguments[1], e.Connection.GetNick()+":"):
				direct = true
			case strings.HasPrefix(e.Arguments[1], e.Connection.GetNick()+","):
				direct = true
			}
		} else {
			private = true
		}

		// if handler is configured to handle the event, then pass it on untouched
		if (hDirect && direct) || (hIndirect && indirect) || (hPrivate && private) {
			f(e)
		}
	}
}
