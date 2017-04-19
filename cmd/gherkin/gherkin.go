package main

import (
	"crypto/tls"
	"flag"
	"strings"

	"github.com/peted27/gherkin/pkg/gherkin"
	"github.com/peted27/gherkin/pkg/plugins/autoban"
	"github.com/peted27/gherkin/pkg/plugins/pinger"
	"github.com/peted27/gherkin/pkg/plugins/sed"
	"github.com/peted27/gherkin/pkg/plugins/seen"
	"github.com/peted27/gherkin/pkg/plugins/slap"
	"github.com/peted27/gherkin/pkg/plugins/uptime"
	"github.com/peted27/gherkin/pkg/plugins/urltitle"
	irc "github.com/peted27/go-ircevent"
)

var (
	host     = flag.String("host", "irc.example.com", "Server host[:port]")
	ssl      = flag.Bool("ssl", true, "Enable SSL")
	nick     = flag.String("nick", "goircbot", "Bot nick")
	ident    = flag.String("ident", "goircbot", "Bot ident")
	channels = flag.String("channels", "", "Channels to join (separated by comma)")
	debug    = flag.Bool("debug", false, "Enable debugging output")
	plugins  = map[string]gherkin.Plugin{}
	version  = "0.9.1"
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

	bot.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsCommandMessage(e) {
				return
			}

			if strings.HasPrefix(e.Arguments[1], "!help") {
				for _, p := range plugins {
					e.Connection.Privmsg(e.Nick, gherkin.MakeHelpString(p))
				}
			}
		})

	bot.AddCallback("PRIVMSG",
		func(e *irc.Event) {
			if !gherkin.IsCommandMessage(e) {
				return
			}

			if strings.HasPrefix(e.Arguments[1], "!version") {
				e.Connection.Action(e.Arguments[0], "running version "+version)
			}
		})

	// plugin registration
	uptime.Register(bot, plugins)
	pinger.Register(bot, plugins)
	slap.Register(bot, plugins)
	urltitle.Register(bot, plugins)
	sed.Register(bot, plugins)
	seen.Register(bot, plugins)
	autoban.Register(bot, plugins)

	bot.Loop()

}
