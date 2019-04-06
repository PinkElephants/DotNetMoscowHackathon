package calc

import (
	"github.com/PinkElephants/DotNetMoscowHackathon/client"
)

const (
	NotBad        = "NotBad"
	Drifted       = "Drifted"
	Hungry        = "Hungry"
	Punished      = "Punished"
	HappyAsInsane = "HappyAsInsane"
)

const (
	East      = "East"
	West      = "West"
	NorthEast = "NorthEast"
	NorthWest = "NorthWest"
	SouthEast = "SouthEast"
	SouthWest = "SouthWest"
)

type Bot struct {
	Help client.Help
	Info client.ServerInfo

	turn client.Turn
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) Result(result client.TurnResult) {

	b.turn = client.Turn{
		Direction:    "West",
		Acceleration: 1,
	}
}

func (b *Bot) Turn() client.Turn {
	return b.turn
}
