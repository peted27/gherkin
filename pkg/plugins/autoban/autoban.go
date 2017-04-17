package autoban

import (
	"sync"
	"time"

	"github.com/peted27/gherkin/pkg/gherkin"
	"github.com/peted27/go-ircevent"
)

var (
	db            = Log{}
	con           *irc.Connection
	timeConnected time.Time
)

// Log is an accessible map of channels to nick to entries.
type Log struct {
	sync.Mutex
	M map[string]map[string]time.Time
}

func Register(c *irc.Connection) {
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
			for channel, chanMap := range db.M {
				for nick, t := range chanMap {
					if time.Since(t).Minutes() >= 1 {
						if con.Debug {
							con.Log.Println("plugin (autoban): user (" + nick + ") banned from " + channel)
						}
						db.Remove(channel, nick)
						c.Kick(nick, channel, "Too late!")
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

	db.Store(channel, nick, time.Now())
	con.Notice(nick, "Welcome to "+channel+", you have 60 seconds to chat or you will be banned.")
}

func onPrivmsg(e *irc.Event) {
	channel := e.Arguments[0]
	nick := e.Nick

	if _, found := db.Search(channel, nick); found == true {
		if con.Debug {
			con.Log.Println("plugin (autoban): user (" + nick + ") acknowledged, removing from auto ban for " + channel)
		}
		db.Remove(channel, nick)
	}
}

func onMode(e *irc.Event) {
	channel := e.Arguments[0]
	mode := e.Arguments[1]
	nick := e.Arguments[2]

	if mode == "+v" || mode == "+o" {
		if _, found := db.Search(channel, nick); found == true {
			if con.Debug {
				con.Log.Println("plugin (autoban): user (" + nick + ") acknowledged, removing from auto ban for " + channel)
			}
			db.Remove(channel, nick)
		}
	}
}

// Store saves a line from a channel/nick into backlog.
func (l *Log) Store(channel, nick string, seen time.Time) {
	l.Lock()
	defer l.Unlock()

	if l.M == nil {
		if con.Debug {
			con.Log.Println("plugin (autoban): creating channel map for " + channel)
		}
		l.M = map[string]map[string]time.Time{}
	}

	if _, p := l.M[channel]; p {
		// update time
		if con.Debug {
			con.Log.Println("plugin (autoban): updating seen time for " + nick)
		}
		l.M[channel][nick] = seen
	} else {
		if con.Debug {
			con.Log.Println("plugin (autoban): creating nick map for " + nick)
		}
		l.M[channel] = map[string]time.Time{nick: seen}
	}
}

// Search returns backlog lines of a channel/nick.
func (l *Log) Search(channel, nick string) (time.Time, bool) {
	var results time.Time
	l.Lock()
	defer l.Unlock()
	if _, p := l.M[channel]; p {
		if _, q := l.M[channel][nick]; q {
			results := l.M[channel][nick]
			if con.Debug {
				con.Log.Println("plugin (autoban): found result")
			}
			return results, true
		}
	}
	if con.Debug {
		con.Log.Println("plugin (autoban): result not found")
	}
	return results, false
}

// Remove deletes an entry from backlog lines of a channel/nick.
func (l *Log) Remove(channel, nick string) bool {

	l.Lock()
	defer l.Unlock()
	if _, p := l.M[channel]; p {
		if _, q := l.M[channel][nick]; q {
			delete(l.M[channel], nick)
			if con.Debug {
				con.Log.Println("plugin (autoban): removing user from db")
			}
			return true
		}
	}
	if con.Debug {
		con.Log.Println("plugin (autoban): user not found")
	}
	return false
}
