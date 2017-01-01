package urltitle

import (
	"errors"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/peted27/gherkin/lib"
	"github.com/peted27/go-ircevent"
)

var (
	linkRE       = regexp.MustCompile(`(?:^|\s)(https?://[^#\s]+)`)
	silenceRE    = regexp.MustCompile(`(^|\s)tg(\)|\s|$)`) // Line ignored if matched.
	titleRE      = regexp.MustCompile(`(?i)<title[^>]*>([^<]+)<`)
	whitespaceRE = regexp.MustCompile(`\s+`)
)

func Register(c *irc.Connection) {
	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !lib.IsChatMessage(e) {
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
		log.Println("urltitle:", err)
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
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	res.Body.Close()

	title, err := Parse(string(body[:]))
	if err != nil {
		return "", err
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
