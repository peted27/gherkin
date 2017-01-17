package slap

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/peted27/gherkin/lib"
	"github.com/peted27/go-ircevent"
)

var (
	slaps   []string
	command = "!slap"
)

func Register(c *irc.Connection) {
	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !lib.IsCommandMessage(e) {
				return
			}
			handle(e)
		})
}

func Initialise() {
	// open a file
	if file, err := os.Open("plugins/slap/slap.txt"); err == nil {

		// make sure it gets closed
		defer file.Close()

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			slaps = append(slaps, scanner.Text())
		}

		// check for errors
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal(err)
	}

	rand.Seed(42)
}

func handle(e *irc.Event) {
	channel := e.Arguments[0]
	text := e.Arguments[1]
	nick := e.Nick
	target := nick

	if !strings.HasPrefix(e.Arguments[1], command) {
		return
	}

	if text != command {
		target = strings.TrimPrefix(text, command+" ")
	}

	e.Connection.Action(channel, "slaps "+target+" around a bit "+slaps[rand.Intn(len(slaps))])
}
