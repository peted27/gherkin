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
	con           *irc.Connection
	timeConnected time.Time

	info = gherkin.Plugin{
		Name:    "seen",
		Command: "!seen <user>",
		Help:    "last time <user> was seen",
		Version: "0.1.0",
	}
)

func Register(c *irc.Connection, h map[string]gherkin.Plugin) {
	h[info.Name] = info
	con = c
	timeConnected = time.Now()

	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsPublicMessage(e) && !gherkin.IsCommandMessage(e) {
				return
			}
			handle(e)
		})
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

	if !strings.HasPrefix(e.Arguments[1], info.Command) {
		return
	}

	if text != info.Command {
		target = strings.TrimPrefix(text, info.Command+" ")
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
