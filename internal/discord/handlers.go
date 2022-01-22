package discord

import (
	"fmt"
	"strings"
	"wiseman/internal/db"

	"github.com/bwmarrin/discordgo"
)

type CommandFunc func(*discordgo.Session, *discordgo.MessageCreate, []string) error

var Commands map[string]CommandFunc

func init() {
	Commands = make(map[string]CommandFunc, 200)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID || len(m.Content) < 1 {
		return
	}

	// Check if prefix for this server is correct
	if db.GetServerByID(m.GuildID).ServerPrefix != m.Content[0:1] {
		return
	}

	msg := strings.Split(m.Content[1:], " ")

	command := strings.ToLower(msg[0])

	args := msg[1:]

	// Check if command exists
	if _, ok := Commands[command]; !ok {
		return
	}

	err := Commands[command](s, m, args)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func memberAdd(s *discordgo.Session, u *discordgo.GuildMemberAdd) {
	fmt.Println("New Member", u.User.Username)
}

func memberRemove(s *discordgo.Session, u *discordgo.GuildMemberRemove) {
	fmt.Println("Member Removed", u.User.Username)
}

func serverAdd(s *discordgo.Session, g *discordgo.GuildCreate) {
	fmt.Println("New Server", g.Name)
}

func serverRemove(s *discordgo.Session, g *discordgo.GuildDelete) {
	fmt.Println("Server Removed", g.ID)
}
