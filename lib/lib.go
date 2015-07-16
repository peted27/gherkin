package lib

import (
	"strings"

	"github.com/peted27/go-ircevent"
)

type Command struct {
	Private bool
	Direct  bool
	Command bool
	Line    bool
}

func PrivmsgHandler(f func(e *irc.Event), c *Command) func(e *irc.Event) {

	return func(e *irc.Event) {
		var private, direct, command, line bool

		// deterimine the type of private message
		if strings.HasPrefix(e.Arguments[0], "#") {
			switch {
			case strings.HasPrefix(e.Arguments[1], "!"):
				command = true
			case strings.HasPrefix(e.Arguments[1], e.Connection.GetNick()+":"):
				direct = true
			case strings.HasPrefix(e.Arguments[1], e.Connection.GetNick()+","):
				direct = true
			default:
				line = true
			}
		} else {
			private = true
		}

		// if handler is configured to handle the event, then pass it on untouched
		if (c.Direct && direct) || (c.Command && command) || (c.Private && private) || (c.Line && line) {
			f(e)
		}
	}
}
