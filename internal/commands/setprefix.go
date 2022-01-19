package commands

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "setPrefix",
		Category:    "This is a category",
		Description: "This is a descriptio",
		Usage:       "This is a usage",
	})
}

func Setprefix() {

}
