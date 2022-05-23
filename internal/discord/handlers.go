package discord

import (
	"fmt"
	"strings"
	"time"
	"wiseman/internal/db"
	"wiseman/internal/services/user"

	"github.com/bwmarrin/discordgo"
)

type userTimer struct {
	UserId  string
	GuildId string
}

type CommandFunc func(*discordgo.Session, *discordgo.MessageCreate, []string) error

var Commands map[string]CommandFunc

// set a ticker to count, second by second,
//the amount of time the user is in the voice channel
var tick = time.NewTicker(1000 * time.Millisecond)

// first element of the array is the joinCh and second is leaveCh
var leaveMap = make(map[string]chan bool)
var changeMap = make(map[string]chan bool)
var joinCh = make(chan userTimer)

func init() {
	Commands = make(map[string]CommandFunc, 200)
	go handleTimers()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID || len(m.Content) < 1 {
		return
	}

	fmt.Println("Message:", m.Content, "from:", m.Author.Username, "in:", m.GuildID)

	u := db.GetUserByID(m.Author.ID, m.GuildID)

	user.IncreaseExperience(u, 10, m.GuildID)
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
	server := db.GetServerByID(u.GuildID)
	s.ChannelMessageSend(server.WelcomeChannel, strings.ReplaceAll(server.WelcomeMessage, "[user]", u.User.Username))
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

// TODO: If the user is in a voice channel while the bot is starting,
// the bot will not be able to track the user, and he needs to reenter the channel
// to be tracked.

// Being multiple people in a channel gives you more points
// Muting yourself gives you less points
// If you never talk you get less points
func voiceStateChange(s *discordgo.Session, c *discordgo.VoiceStateUpdate) {

	// Leave the voice channel
	if c.BeforeUpdate != nil {
		// check if the user is streaming his screen
		if c.ChannelID == "" {
			fmt.Println(c.VoiceState.UserID, "left", c.GuildID)
			leaveMap[c.UserID] <- true
		} else if c.ChannelID != "" && c.ChannelID != c.BeforeUpdate.ChannelID {
			fmt.Println("The user changed voice channel")
			changeMap[c.UserID] <- true
			joinCh <- userTimer{
				UserId:  c.UserID,
				GuildId: c.GuildID,
			}
		} else if c.ChannelID != "" && !c.BeforeUpdate.SelfMute && c.SelfMute && !c.SelfDeaf {
			fmt.Println("The user muted himself")
		} else if c.ChannelID != "" && !c.BeforeUpdate.SelfDeaf && c.SelfDeaf {
			fmt.Println("The user deafened himself")
		} else if c.ChannelID != "" && c.BeforeUpdate.SelfMute && !c.BeforeUpdate.SelfDeaf && !c.SelfMute && !c.SelfDeaf {
			fmt.Println("The user unmuted himself")
		} else if c.ChannelID != "" && c.BeforeUpdate.SelfDeaf && c.BeforeUpdate.SelfMute && !c.SelfDeaf && !c.SelfMute {
			fmt.Println("The user undeafened himself")
		} else {
			fmt.Println("The user is performing another operation")
		}
	} else {
		// Join
		fmt.Println(c.UserID, "Joined", c.GuildID)
		joinCh <- userTimer{
			UserId:  c.UserID,
			GuildId: c.GuildID,
		}
	}
}

func handleTimers() {

	for {

		now := time.Now().Unix()
		counter := 0
		e := <-joinCh

		go func(ut userTimer) {
			changeMap[ut.UserId] = make(chan bool)

			for {
				select {
				case <-leaveMap[ut.UserId]:
					fmt.Println("Time Spent", counter, "seconds; from", time.Unix(now, 0), "to", time.Unix(now+int64(counter), 0))
					close(leaveMap[ut.UserId])
					delete(leaveMap, ut.UserId)
					return
				case <-changeMap[ut.UserId]:
					fmt.Println("Time Spent", counter, "seconds; from", time.Unix(now, 0), "to", time.Unix(now+int64(counter), 0))
					counter = 0
					now = time.Now().Unix()
				case <-tick.C:
					counter += 1
					u := db.GetUserByID(ut.UserId, ut.GuildId)
					user.IncreaseExperience(u, uint(time.Second*1), ut.GuildId)
				}
			}
		}(e)
	}
}
