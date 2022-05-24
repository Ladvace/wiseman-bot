package entities

type UserType struct {
	ComplexID              string `bson:"complexid"`
	ServerID               string `bson:"serverid"`
	UserID                 string `bson:"userid"`
	MessagesCount          uint   `bson:"messagescount"`
	CurrentLevelExperience uint   `bson:"currentlevelexperience"`
	CurrentLevel           uint   `bson:"currentlevel"`
	LastTimeOnline         uint64 `bson:"lastranktime"`
	Bot                    bool   `bson:"bot"`
	Verified               bool   `bson:"verified"`
}

type UsersType map[string]*UserType
