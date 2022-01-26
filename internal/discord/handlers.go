package discord

import (
	"fmt"
	"strings"
	"time"
	"wiseman/internal/db"

	"github.com/bwmarrin/discordgo"
)

type CommandFunc func(*discordgo.Session, *discordgo.MessageCreate, []string) error

var Commands map[string]CommandFunc

var joinTimestamps map[string]int64

func init() {
	Commands = make(map[string]CommandFunc, 200)
	joinTimestamps = make(map[string]int64, 1000)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID || len(m.Content) < 1 {
		return
	}

	server := db.GetServerByID(m.GuildID)

	fmt.Println("Message:", m.Content)

	u := db.GetUserByID(m.Author.ID, m.GuildID)
	u.IncreaseExperience(10, uint(server.MsgExpMultiplier))
	fmt.Println("After Message:", m.Content)

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

func memberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	fmt.Println("Member Updated", m.User.Username)
}

func serverAdd(s *discordgo.Session, g *discordgo.GuildCreate) {
	fmt.Println("New Server", g.Name)
}

func serverRemove(s *discordgo.Session, g *discordgo.GuildDelete) {
	fmt.Println("Server Removed", g.ID)
}

// Being multiple people in a channel gives you more points
// Muting yourself gives you less points
// If you never talk you get less points
func voiceStateChange(s *discordgo.Session, c *discordgo.VoiceStateUpdate) {
	if c.BeforeUpdate != nil {
		server := db.GetServerByID(c.GuildID)

		evStr := fmt.Sprintf("%s %s %s", c.GuildID, c.BeforeUpdate.ChannelID, c.UserID)
		_, ok := joinTimestamps[evStr]
		if !ok {
			return
		}
		timeDiff := time.Now().Unix() - joinTimestamps[evStr]
		// Leave
		fmt.Println("Left after", timeDiff, "seconds")
		delete(joinTimestamps, evStr)
		u := db.GetUserByID(c.UserID, c.GuildID)
		u.IncreaseExperience(uint(timeDiff)*2, uint(server.TimeExpMultiplier))
	} else {
		evStr := fmt.Sprintf("%s %s %s", c.GuildID, c.ChannelID, c.UserID)
		// Join
		fmt.Println("Joined")
		joinTimestamps[evStr] = time.Now().Unix()
	}
}
