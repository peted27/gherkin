package main

import (
	"crypto/tls"
	"flag"
	"strings"

	"github.com/peted27/gherkin/lib"
	"github.com/peted27/gherkin/plugins/urltitle"
	"github.com/peted27/go-ircevent"
)

var (
	host     = flag.String("host", "irc.example.com", "Server host[:port]")
	ssl      = flag.Bool("ssl", true, "Enable SSL")
	nick     = flag.String("nick", "goircbot", "Bot nick")
	ident    = flag.String("ident", "goircbot", "Bot ident")
	channels = flag.String("channels", "", "Channels to join (separated by comma)")
	debug    = flag.Bool("debug", false, "Enable debugging output")
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
	bot.AddCallback("PRIVMSG", lib.PrivmsgHandler(
		func(e *irc.Event) {
			if strings.HasPrefix(e.Arguments[1], "!ping") {
				e.Connection.Privmsg(e.Arguments[0], "pong!")
			}
		}, false, true, false))

	// slap plugin
	bot.AddCallback("PRIVMSG", lib.PrivmsgHandler(
		func(e *irc.Event) {
			if strings.HasPrefix(e.Arguments[1], "!slap") {
				e.Connection.Action(e.Arguments[0], "slaps "+e.Nick+" around a bit with a large trout!")
			}
		}, false, true, false))

	urltitle.Register(bot)

	bot.Loop()

}
