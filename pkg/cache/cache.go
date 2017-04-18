package cache

import "sync"

// Cache is an accessible map of channels to nick to entries.
type Cache struct {
	sync.Mutex
	M map[string]map[string]Storable
}

// Storable in to cache
type Storable struct {
	Field interface{}
}

// Store saves a line from a channel/nick
func (l *Cache) Store(channel, nick string, obj Storable) {
	l.Lock()
	defer l.Unlock()

	if l.M == nil {
		l.M = map[string]map[string]Storable{}
	}

	if _, p := l.M[channel]; p {
		l.M[channel][nick] = obj
	} else {
		l.M[channel] = map[string]Storable{nick: obj}
	}
}

// Search returns backlog lines of a channel/nick.
func (l *Cache) Search(channel, nick string) (Storable, bool) {
	var results Storable
	l.Lock()
	defer l.Unlock()
	if _, p := l.M[channel]; p {
		if _, q := l.M[channel][nick]; q {
			results := l.M[channel][nick]

			return results, true
		}
	}

	return results, false
}

// Remove deletes an entry
func (l *Cache) Remove(channel, nick string) bool {

	l.Lock()
	defer l.Unlock()
	if _, p := l.M[channel]; p {
		if _, q := l.M[channel][nick]; q {
			delete(l.M[channel], nick)
			return true
		}
	}

	return false
}
