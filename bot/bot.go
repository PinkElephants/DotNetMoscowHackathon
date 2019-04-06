package calc

import (
	"github.com/PinkElephants/DotNetMoscowHackathon/client"
)

type Bot struct {
	Help client.Help
}

func NewBot() *Bot {
	return &Bot{}
}
