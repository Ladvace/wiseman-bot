package entities

type RankType struct {
	RankName     string `bson:"rankname"`
	RankMinLevel uint   `bson:"rankminlevel"`
}

type ServerType struct {
	ServerID            string     `bson:"serverid"`
	ServerPrefix        string     `bson:"guildprefix"`
	NotificationChannel string     `bson:"notificationchannel"`
	WelcomeChannel      string     `bson:"welcomechannel"`
	CustomRanks         []RankType `bson:"customranks"`
	RankTime            int        `bson:"ranktime"`
	MsgExpMultiplier    float64    `bson:"msgexpmultiplier"`
	TimeExpMultiplier   float64    `bson:"timeexpmultiplier"`
	WelcomeMessage      string     `bson:"welcomemessage"`
	DefaultRole         string     `bson:"defaultrole"`
}

type ServersType map[string]ServerType
