package seen

import (
	"strings"
	"time"

	"github.com/peted27/gherkin/pkg/cache"
	"github.com/peted27/gherkin/pkg/gherkin"
	"github.com/peted27/go-ircevent"
)

var (
	db            = cache.Cache{}
	command       = "!seen"
	help          = "last time <user> was seen"
	con           *irc.Connection
	timeConnected time.Time
)

func Register(c *irc.Connection, h map[string]string) {
	con = c
	timeConnected = time.Now()
	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsPublicMessage(e) && !gherkin.IsCommandMessage(e) {
				return
			}
			handle(e)
		})
	h[command] = help
}

func handle(e *irc.Event) {
	channel := e.Arguments[0]
	text := e.Arguments[1]
	nick := e.Nick
	target := nick

	db.Store(channel, nick, cache.Storable{Field: time.Now()})
	if con.Debug {
		con.Log.Println("plugin (seen): user (" + nick + ") updating time on " + channel)
	}

	if !strings.HasPrefix(e.Arguments[1], command) {
		return
	}

	if text != command {
		target = strings.TrimPrefix(text, command+" ")
		target = strings.TrimSpace(target)
	}

	if t, found := db.Search(channel, target); found == false {
		if con.Debug {
			con.Log.Println("plugin (seen): user (" + nick + ") not seen on " + channel)
		}
		e.Connection.Action(channel, target+" has not been seen (bot online since "+timeConnected.Format("15:04:05 (2006-01-02) MST")+")")
	} else {
		if con.Debug {
			con.Log.Println("plugin (seen): user (" + nick + ") seen on " + channel)
		}
		e.Connection.Action(channel, target+" last seen "+t.Field.(time.Time).Format("15:04:05 (2006-01-02) MST"))
	}
}
