package commands

import "fmt"

type Helper struct {
	Name        string
	Category    string
	Description string
	Usage       string
}

var Helpers []Helper

func Help() {
	for _, v := range Helpers {
		fmt.Println("Helper", v.Description)
	}
}
