package uptime

import (
	"strings"
	"time"

	"github.com/peted27/gherkin/pkg/gherkin"
	irc "github.com/peted27/go-ircevent"
)

var (
	timeInitialised time.Time
	con             *irc.Connection
	info            = gherkin.Plugin{
		Name:    "uptime",
		Command: "!uptime",
		Help:    "display time since bot was launched",
		Version: "0.1.0",
	}
)

func Register(c *irc.Connection, h map[string]gherkin.Plugin) {
	h[info.Name] = info
	con = c
	timeInitialised = time.Now()

	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsCommandMessage(e) {
				return
			}

			if strings.HasPrefix(e.Arguments[1], "!uptime") {
				e.Connection.Action(e.Arguments[0], "running since "+timeInitialised.Format("15:04:05 (2006-01-02) MST"))
			}
		})
}
