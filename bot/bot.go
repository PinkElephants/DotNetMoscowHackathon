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

	info client.ServerInfo

	turn  client.Turn
	car   client.Car
	cells []client.Cell
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) Result(result client.TurnResult) {
	b.car = result.Car()
	b.cells = result.Cells()

	b.turn = client.Turn{
		Direction:    "West",
		Acceleration: b.acceleration(),
	}
}

func (b *Bot) Start(info client.ServerInfo) {
	b.info = info
	b.car = info.Car()
	b.cells = info.Cells()
}

func (b *Bot) Turn() client.Turn {
	return b.turn
}

func (b *Bot) acceleration() int {
	safeSpeed := b.Help.MaxSpeed
	for _, a := range b.Help.DriftsAngles {
		if safeSpeed > a.MaxSpeed {
			safeSpeed = a.MaxSpeed
		}
	}
	return safeSpeed - b.car.Speed
}
