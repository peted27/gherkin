package slap

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/peted27/gherkin/pkg/gherkin"
	"github.com/peted27/go-ircevent"
)

var (
	con   *irc.Connection
	slaps []string

	info = gherkin.Plugin{
		Name:    "slap",
		Command: "!slap <user>",
		Help:    "randomly slap <user>",
		Version: "0.1.0",
	}
)

func Register(c *irc.Connection, h map[string]gherkin.Plugin) {
	h[info.Name] = info
	con = c
	initialise()

	c.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsCommandMessage(e) {
				return
			}
			handle(e)
		})
}

func initialise() {
	// open a file
	if file, err := os.Open("pkg/plugins/slap/slap.txt"); err == nil {

		// make sure it gets closed
		defer file.Close()

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			slaps = append(slaps, scanner.Text())

		}
		if con.Debug {
			con.Log.Println("plugin (slap): loaded slaps from file")
		}

		// check for errors
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal(err)
	}

	// random number seed, order can be played with by changing this seed
	rand.Seed(42)
}

func handle(e *irc.Event) {
	channel := e.Arguments[0]
	text := e.Arguments[1]
	nick := e.Nick
	target := nick

	if !strings.HasPrefix(e.Arguments[1], info.Command) {
		return
	}

	if text != info.Command {
		target = strings.TrimPrefix(text, info.Command+" ")
		target = strings.TrimSpace(target)
	}

	e.Connection.Action(channel, "slaps "+target+" around a bit "+slaps[rand.Intn(len(slaps))])
}
