package pinger

import (
	"strings"

	"github.com/peted27/gherkin/pkg/gherkin"
	irc "github.com/peted27/go-ircevent"
)

var (
	con  *irc.Connection
	info = gherkin.Plugin{
		Name:    "ping",
		Command: "!ping",
		Help:    "reply to user with !pong",
		Version: "0.1.0",
	}
)

func Register(c *irc.Connection, h map[string]gherkin.Plugin) {
	h[info.Name] = info
	con = c

	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsCommandMessage(e) {
				return
			}

			if strings.HasPrefix(e.Arguments[1], "!ping") {
				e.Connection.Privmsg(e.Arguments[0], "pong!")
			}
		})
}
