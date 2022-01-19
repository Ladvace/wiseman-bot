package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Helper struct {
	Name        string
	Category    string
	Description string
	Usage       string
}

var Helpers []Helper

func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, v := range Helpers {
		fmt.Println("Helper", v.Description)
	}
}
