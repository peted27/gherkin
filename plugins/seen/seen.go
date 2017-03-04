package seen

import (
	"strings"
	"sync"
	"time"

	"github.com/peted27/gherkin/lib"
	"github.com/peted27/go-ircevent"
)

var (
	db            = Log{}
	command       = "!seen"
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
	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !lib.IsPublicMessage(e) && !lib.IsCommandMessage(e) {
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

	db.Store(channel, nick, time.Now())

	if !strings.HasPrefix(e.Arguments[1], command) {
		return
	}

	if text != command {
		target = strings.TrimPrefix(text, command+" ")
		target = strings.TrimSpace(target)
	}

	if t, found := db.Search(channel, target); found == false {
		e.Connection.Action(channel, target+" has not been seen (bot online since "+timeConnected.Format("15:04:05 (2006-01-02) MST")+")")
	} else {
		e.Connection.Action(channel, target+" last seen "+t.Format("15:04:05 (2006-01-02) MST"))
	}
}

// Store saves a line from a channel/nick into backlog.
func (l *Log) Store(channel, nick string, seen time.Time) {
	l.Lock()
	defer l.Unlock()

	if l.M == nil {
		if con.Debug {
			con.Log.Println("plugin (seen): creating channel map for " + channel)
		}
		l.M = map[string]map[string]time.Time{}
	}

	if _, p := l.M[channel]; p {
		// update time
		if con.Debug {
			con.Log.Println("plugin (seen): updating seen time for " + nick)
		}
		l.M[channel][nick] = seen
	} else {
		if con.Debug {
			con.Log.Println("plugin (seen): creating nick map for " + nick)
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
				con.Log.Println("plugin (seen): found result")
			}
			return results, true
		}
	}
	if con.Debug {
		con.Log.Println("plugin (seen): result not found")
	}
	return results, false
}
