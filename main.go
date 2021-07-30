package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var Log *log.Logger

func init() {

}

type NattyTracker struct {
	User       string
	NattyCount int
}

func main() {
	Log = log.Default()
	discord, err := discordgo.New("Bot " + "<token>")
	if err != nil {
		Log.Println("error creating Discord session,", err)
		return
	}

	if _, err := os.Stat("stats.json"); os.IsNotExist(err) {
		if err != nil {
			Log.Print(err)
		}
		Log.Println("Initializing stats.json")
		os.Create("stats.json")
	} else {
		Log.Println("stats.json exists")
		stats, err := os.Open("stats.json")
		if err != nil {
			Log.Print(err)
		}
		Log.Println("Opened stats.json")
		defer stats.Close()
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		Log.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "natty") {
		Log.Printf("Message Channel ID is %s, Message ID is %s", m.ChannelID, m.ID)
		s.ChannelMessageSend(m.ChannelID, "HELL YEAH BROTHER")
		s.MessageReactionAdd(m.ChannelID, m.ID, ":nattyd:867961442039648338")
		Log.Printf("Sent some natty daddy encouragement to %s", m.Author.Username)
		data := NattyTracker{
			User:       m.Author.Username,
			NattyCount: 1,
		}
		file, _ := json.MarshalIndent(data, "", " ")

		_ = ioutil.WriteFile("stats.json", file, 0644)
	}

}
