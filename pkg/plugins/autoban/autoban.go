package autoban

import (
	"strings"
	"time"

	"github.com/peted27/gherkin/pkg/cache"
	"github.com/peted27/gherkin/pkg/gherkin"
	"github.com/peted27/go-ircevent"
)

var (
	joined        = cache.Cache{}
	banned        = cache.Cache{}
	con           *irc.Connection
	timeConnected time.Time
	info          = gherkin.Plugin{
		Name:    "autoban",
		Command: "",
		Help:    "autobans bots from channel",
		Version: "0.1.0",
	}
)

func Register(c *irc.Connection, h map[string]gherkin.Plugin) {
	h[info.Name] = info
	con = c
	timeConnected = time.Now()
	c.AddCallback("JOIN",
		func(e *irc.Event) {
			onJoin(e)
		})
	c.AddCallback("MODE",
		func(e *irc.Event) {
			onMode(e)
		})
	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsPublicMessage(e) && !gherkin.IsCommandMessage(e) {
				return
			}
			onPrivmsg(e)
		})

	go func() {
		for {
			time.Sleep(60 * time.Second)
			// kick and band all users who havnt spoken
			for channel, chanMap := range joined.M {
				for nick, t := range chanMap {
					if time.Since(t.Field.(time.Time)).Minutes() >= 1 {
						if con.Debug {
							con.Log.Println("plugin (autoban): user (" + nick + ") banned from " + channel)
						}
						joined.Remove(channel, nick)
						banned.Store(channel, nick, cache.Storable{Field: time.Now()})
						// naiive ban just on nick for now
						c.Mode(channel, "+b "+nick)
						c.Kick(nick, channel, "Too late!")
					}
				}
			}

			//unban users after 60 minutes
			for channel, chanMap := range banned.M {
				for nick, t := range chanMap {
					if time.Since(t.Field.(time.Time)).Minutes() >= 60 {
						if con.Debug {
							con.Log.Println("plugin (autoban): user (" + nick + ") un-banned from " + channel)
						}
						banned.Remove(channel, nick)
						c.Mode(channel, "-b "+nick)
					}
				}
			}
		}
	}()

}

func onJoin(e *irc.Event) {
	channel := e.Arguments[0]
	nick := e.Nick

	if nick == con.GetNick() {
		return
	}

	if con.Debug {
		con.Log.Println("plugin (autoban): user (" + nick + ") joined, starting auto ban timer for " + channel)
	}

	joined.Store(channel, nick, cache.Storable{Field: time.Now()})
	con.Notice(nick, "Welcome to "+channel+", you have 60 seconds to chat or you will be banned.")
}

func onPrivmsg(e *irc.Event) {
	channel := e.Arguments[0]
	nick := e.Nick

	if _, found := joined.Search(channel, nick); found == true {
		if con.Debug {
			con.Log.Println("plugin (autoban): user (" + nick + ") acknowledged, removing from auto ban for " + channel)
		}
		joined.Remove(channel, nick)
	}
}

func onMode(e *irc.Event) {
	channel := e.Arguments[0]

	if !strings.HasPrefix(channel, "#") {
		return
	}
	// just in case not a channel
	mode := e.Arguments[1]
	nick := e.Arguments[2]

	if mode == "+v" || mode == "+o" {
		if _, found := joined.Search(channel, nick); found == true {
			if con.Debug {
				con.Log.Println("plugin (autoban): user (" + nick + ") acknowledged, removing from auto ban for " + channel)
			}
			joined.Remove(channel, nick)
		}
	}
}
