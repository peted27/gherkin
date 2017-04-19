package urltitle

import (
	"errors"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/peted27/gherkin/pkg/gherkin"
	"github.com/peted27/go-ircevent"
)

var (
	con          *irc.Connection
	linkRE       = regexp.MustCompile(`(?:^|\s)(https?://[^#\s]+)`)
	silenceRE    = regexp.MustCompile(`(^|\s)tg(\)|\s|$)`) // Line ignored if matched.
	titleRE      = regexp.MustCompile(`(?i)<title[^>]*>([^<]+)<`)
	whitespaceRE = regexp.MustCompile(`\s+`)

	info = gherkin.Plugin{
		Name:    "urltitle",
		Command: "",
		Help:    "automatically retrieves url titles",
		Version: "0.1.0",
	}
)

func Register(c *irc.Connection, h map[string]gherkin.Plugin) {
	h[info.Name] = info
	con = c

	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsPublicMessage(e) {
				return
			}
			handle(e)
		})
}

func handle(e *irc.Event) {
	target := e.Arguments[0]
	text := e.Arguments[1]

	if silenceRE.MatchString(text) {
		return
	}

	match := linkRE.FindStringSubmatch(text)
	if match == nil {
		return
	}

	link := match[1]
	if len(link) > 200 {
		return
	}

	title, err := GetTitle(link)
	if err != nil {
		if con.Debug {
			con.Log.Println("plugin (urltitle): error ", err)
		}
		return
	}

	if len(title) > 200 {
		title = title[:200]
	}

	//e.Connection.Privmsgf(target, "%s :: %s", link, title)
	e.Connection.Privmsgf(target, "<title> :: %s", title)
}

func GetTitle(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		if con.Debug {
			con.Log.Println("plugin (urltitle): could not retrieve url ", err)
		}
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if con.Debug {
			con.Log.Println("plugin (urltitle): could not read body of response ", err)
		}
		return "", err
	}

	res.Body.Close()

	title, err := Parse(string(body[:]))
	if err != nil {
		if con.Debug {
			con.Log.Println("plugin (urltitle): couldnt not parse title ", err)
		}
		return "", err
	}
	if con.Debug {
		con.Log.Println("plugin (urltitle): found title " + title)
	}
	return title, nil
}

func Parse(body string) (string, error) {
	text := titleRE.FindStringSubmatch(body)
	if text == nil {
		return "", errors.New("url: cannot parse title")
	}
	return Trim(html.UnescapeString(text[1])), nil
}

// Trim removes all white spaces including duplicates in a string.
func Trim(s string) string {
	return strings.TrimSpace(whitespaceRE.ReplaceAllString(s, " "))
}
