package sed

import (
	"container/list"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/peted27/gherkin/lib"
	"github.com/peted27/go-ircevent"
)

var (
	con        *irc.Connection
	bl         = Backlog{}
	expiration = 5 * time.Minute
	maxLines   = 4
)

func Register(c *irc.Connection) {
	con = c
	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !lib.IsPublicMessage(e) {
				return
			}
			handle(e)
		})
}

func handle(e *irc.Event) {
	channel := e.Arguments[0]
	text := e.Arguments[1]
	nick := e.Nick

	if !strings.HasPrefix(channel, "#") {
		return
	}

	r, err := regexp.Compile("^s/([^/]+)/([^/]*)/g?")
	if err != nil {
		e.Connection.Privmsg(channel, fmt.Sprintf("error: %s", err))
		return
	}

	m := r.FindSubmatch([]byte(text))
	if m == nil {
		bl.Store(channel, nick, text)
		return
	}

	meant := bl.Sed(channel, nick, string(m[1]), string(m[2]))
	if meant == "" {
		return
	}

	e.Connection.Privmsg(channel, fmt.Sprintf("%s meant: %s", nick, meant))
	bl.Store(channel, nick, meant)

}

// Entry is a backlog entry: text line and time it happened.
type Entry struct {
	Line string
	Time time.Time
}

// Expired tells if an entry is expired based on time it happened and expiration.
func (e *Entry) Expired() bool {
	return time.Now().After(e.Time.Add(expiration))
}

// Backlog is an accessible map of channels to nick to entries.
type Backlog struct {
	sync.Mutex
	M map[string]map[string]*list.List
}

// Clean maintains a backlog clean by expiring old entries and ensuring maximum of lines.
func (bl *Backlog) Clean() {
	bl.Lock()
	defer bl.Unlock()
	for channel, cl := range bl.M {
		for nick, l := range cl {
			var rm []*list.Element
			i := 0
			for e := l.Front(); e != nil; e = e.Next() {
				entry := e.Value.(Entry)
				if entry.Expired() || i >= maxLines {
					rm = append(rm, e)
				} else {
					i++
				}
			}
			for _, e := range rm {
				l.Remove(e)
			}
			if l.Len() == 0 {
				delete(cl, nick)
			}
		}
		if len(cl) == 0 {
			delete(bl.M, channel)
		}
	}
}

// Store saves a line from a channel/nick into backlog.
func (bl *Backlog) Store(channel, nick, line string) {
	defer bl.Clean()
	bl.Lock()
	defer bl.Unlock()
	e := Entry{line, time.Now()}
	if bl.M == nil {
		bl.M = map[string]map[string]*list.List{}
	}
	if _, p := bl.M[channel]; p {
		if _, q := bl.M[channel][nick]; q {
			// most recent line is in first position
			bl.M[channel][nick].PushFront(e)
		} else {
			l := list.New()
			l.PushBack(e)
			bl.M[channel][nick] = l
		}
	} else {
		l := list.New()
		l.PushBack(e)
		bl.M[channel] = map[string]*list.List{nick: l}
	}
}

// Search returns backlog lines of a channel/nick.
func (bl *Backlog) Search(channel, nick string) []string {
	var results []string
	bl.Lock()
	defer bl.Unlock()
	if _, p := bl.M[channel]; p {
		if _, q := bl.M[channel][nick]; q {
			l := bl.M[channel][nick]
			for e := l.Front(); e != nil; e = e.Next() {
				entry := e.Value.(Entry)
				results = append(results, entry.Line)
			}
		}
	}
	return results
}

// Sed attempts to replace a pattern in a backlog for channel/nick.
func (bl *Backlog) Sed(channel, nick, pattern, replace string) string {
	if len(pattern) > 80 {
		pattern = pattern[:80]
	}
	if len(replace) > 80 {
		replace = replace[:80]
	}
	for _, line := range bl.Search(channel, nick) {
		if strings.Contains(line, pattern) {
			if con.Debug {
				con.Log.Println("plugin (sed): pattern found, replacing.")
			}
			r := strings.Replace(line, pattern, replace, 1)
			if len(r) > 160 {
				r = r[:160]
			}
			return r
		}
	}
	return ""
}
