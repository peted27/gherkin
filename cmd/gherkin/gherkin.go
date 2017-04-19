package main

import (
	"crypto/tls"
	"flag"
	"strings"
	"time"

	"github.com/peted27/gherkin/pkg/gherkin"
	"github.com/peted27/gherkin/pkg/plugins/autoban"
	"github.com/peted27/gherkin/pkg/plugins/sed"
	"github.com/peted27/gherkin/pkg/plugins/seen"
	"github.com/peted27/gherkin/pkg/plugins/slap"
	"github.com/peted27/gherkin/pkg/plugins/urltitle"
	irc "github.com/peted27/go-ircevent"
)

var (
	host        = flag.String("host", "irc.example.com", "Server host[:port]")
	ssl         = flag.Bool("ssl", true, "Enable SSL")
	nick        = flag.String("nick", "goircbot", "Bot nick")
	ident       = flag.String("ident", "goircbot", "Bot ident")
	channels    = flag.String("channels", "", "Channels to join (separated by comma)")
	debug       = flag.Bool("debug", false, "Enable debugging output")
	helpStrings = map[string]string{}
)

func main() {
	flag.Parse()

	bot := irc.IRC(*nick, *ident)

	//bot.VerboseCallbackHandler = *debug
	bot.Debug = *debug

	// using ssl? configure here
	if *ssl {
		bot.UseTLS = *ssl
		bot.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// connect the bot
	if err := bot.Connect(*host); err != nil {
		bot.Log.Printf("Error: %s\n", err)
	}

	// setup callbacks to join managed channels
	for _, ch := range strings.Split(*channels, ",") {

		bot.AddCallback("001",
			func(e *irc.Event) {
				bot.Join(ch)
				bot.Log.Printf("bot: joining channel %s\n", ch)
			})

	}

	// pong! plugin
	helpStrings["!ping"] = "auto reply with !pong"
	bot.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsCommandMessage(e) {
				return
			}

			if strings.HasPrefix(e.Arguments[1], "!ping") {
				e.Connection.Privmsg(e.Arguments[0], "pong!")
			}
		})

	// !uptime plugin
	timeInitialised := time.Now()
	helpStrings["!uptime"] = "display time since bot was launched"
	bot.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsCommandMessage(e) {
				return
			}

			if strings.HasPrefix(e.Arguments[1], "!uptime") {
				e.Connection.Action(e.Arguments[0], "running since "+timeInitialised.Format("15:04:05 (2006-01-02) MST"))
			}
		})

	// !help
	helpStrings["!help"] = "print this message"
	bot.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsCommandMessage(e) {
				return
			}

			if strings.HasPrefix(e.Arguments[1], "!help") {
				for h, c := range helpStrings {
					e.Connection.Privmsg(e.Nick, gherkin.MakeHelpString(h, c))
				}
			}
		})

	// plugin registration
	slap.Register(bot, helpStrings)
	urltitle.Register(bot, helpStrings)
	sed.Register(bot, helpStrings)
	seen.Register(bot, helpStrings)
	autoban.Register(bot, helpStrings)

	bot.Loop()

}
