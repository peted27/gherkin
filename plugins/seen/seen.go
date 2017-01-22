package seen

import (
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/peted27/gherkin/lib"
	"github.com/peted27/go-ircevent"
)

var (
	db      = Log{}
	command = "!seen"
	con     *irc.Connection
)

// Log is an accessible map of channels to nick to entries.
type Log struct {
	sync.Mutex
	M map[string]map[string]time.Time
}

func Register(c *irc.Connection) {
	con = c
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
	}

	if t, err := db.Search(channel, target); err != nil {
		e.Connection.Action(channel, target+" has not been seen")
	} else {
		e.Connection.Action(channel, target+" last seen "+t.String())
	}
}

// Store saves a line from a channel/nick into backlog.
func (l *Log) Store(channel, nick string, seen time.Time) {
	l.Lock()
	defer l.Unlock()

	if l.M == nil {
		if con.Debug {
			con.Log.Println("plugin (seen): creating channel map")
		}
		l.M = map[string]map[string]time.Time{}
	}

	if _, p := l.M[channel]; p {
		// update time
		if con.Debug {
			con.Log.Println("plugin (seen): updating seen time")
		}
		l.M[channel][nick] = seen
	} else {
		if con.Debug {
			con.Log.Println("plugin (seen): creating nick map")
		}
		l.M[channel] = map[string]time.Time{nick: seen}
	}
}

// Search returns backlog lines of a channel/nick.
func (l *Log) Search(channel, nick string) (time.Time, error) {
	var results time.Time
	l.Lock()
	defer l.Unlock()
	if _, p := l.M[channel]; p {
		if _, q := l.M[channel][nick]; q {
			results := l.M[channel][nick]
			if con.Debug {
				con.Log.Println("plugin (seen): found result")
			}
			return results, nil
		}
	}
	if con.Debug {
		con.Log.Println("plugin (seen): result not found")
	}
	return results, errors.New("Not found")
}
