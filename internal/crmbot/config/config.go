package config

type Config struct {
	MongoDB MongoCfg `json:"mongo"`
	Bot     Bot      `json:"bot"`
}

type Bot struct {
	Token     Token `json:"token" env:"CRM_BOT_THE_SEQUEL_TOKEN,required"`
	CreatorID int64 `json:"creator_id" env:"CREATOR_ID"`
}

type MongoCfg struct {
	Addr     string `json:"addr" env:"DATABASE_URL"`
	HostName string `json:"host_name" env:"MONGO_HOST,required"`
	Port     string `json:"port" env:"MONGO_PORT" default:"27017"`
	DBName   string `json:"db_name" env:"MONGO_DBNAME,required"`
	Username string `json:"username" env:"MONGO_USERNAME,required"`
	Password string `json:"password" env:"MONGO_PASSWORD,required"`
}

func (b Bot) GetToken() Token {
	return b.Token
}

type Token string

func (t Token) String() string {
	return string(t)
}
