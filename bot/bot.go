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

const (
	Empty         = "Empty"
	Rock          = "Rock"
	DangerousArea = "DangerousArea"
	Pit           = "Pit"
)

type Bot struct {
	Help client.Help

	info client.ServerInfo

	turn       client.Turn
	car        client.Car
	cellsIndex [][][]client.Cell
	//cells      []client.Cell
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) Start(info client.ServerInfo) {
	b.info = info
	b.car = info.Car()
	b.allocateCells(info.Cells())
}

func (b *Bot) Result(result client.TurnResult) {
	b.car = result.Car()
	b.allocateCells(result.Cells())
	b.scan()
	// toTarget := b.closestToTarget()

	b.turn = client.Turn{
		Direction:    "West",
		Acceleration: b.acceleration(),
	}
}

func (b *Bot) Turn() client.Turn {
	return b.turn
}

func (b *Bot) allocateCells(cells []client.Cell) {
	b.cellsIndex = make([][][]client.Cell, b.info.Radius)
	for i := range cells {
		b.cellsIndex[i] = make([][]client.Cell, b.info.Radius)
		for j := range b.cellsIndex[i] {
			b.cellsIndex[i][j] = make([]client.Cell, b.info.Radius)
		}
	}
	for _, c := range cells {
		c.Visible = true
		b.cellsIndex[c.X][c.Y][c.Z] = c
	}
}

func (b *Bot) updateCells(cells []client.Cell) {
	// todo
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

func (b *Bot) scan() {

	b.iterAll(func(c client.Cell) {
		if c.Type == "Rock" {
			return
		}

		toCar := c.DistanceFrom(client.Cell{
			X: b.car.X,
			Y: b.car.Y,
			Z: b.car.Z,
		})
		c.DistToCar = toCar
		toTarget := c.DistanceFrom(client.Cell{
			X: b.info.Finish.X,
			Y: b.info.Finish.Y,
			Z: b.info.Finish.Z,
		})
		c.DistToTarget = toTarget
	})
}

func (b *Bot) closestToTarget() client.Cell {

	minDist := 0
	var dist client.Cell

	calcMinDist := func(c client.Cell) {
		if c.Type == Rock {
			return
		}
		if c.DistToTarget < minDist {
			minDist = c.DistToTarget
			dist = c
		}
	}

	b.iterAll(calcMinDist)

	return dist
}

func (b *Bot) iterateNeighbors(cell client.Cell, f func(c client.Cell)) {
	f(b.cellsIndex[cell.X-1][cell.Z][cell.Y+1])
	f(b.cellsIndex[cell.X][cell.Z-1][cell.Y+1])
	f(b.cellsIndex[cell.X+1][cell.Z-1][cell.Y])
	f(b.cellsIndex[cell.X+1][cell.Z][cell.Y-1])
	f(b.cellsIndex[cell.X][cell.Z+1][cell.Y-1])
	f(b.cellsIndex[cell.X-1][cell.Z+1][cell.Y])
}

func (b *Bot) iterAll(f func(c client.Cell)) {
	for x := range b.cellsIndex {
		for y := range b.cellsIndex[x] {
			for z := range b.cellsIndex[x][y] {
				f(b.cellsIndex[x][y][z])
			}
		}
	}
}

func Side(from client.Cell, to client.Cell) string {
	northEast := client.Cell{X: from.X + 1, Y: from.Y - 1, Z: from.Z}
	northWest := client.Cell{X: from.X, Y: from.Y - 1, Z: from.Z + 1}
	west := client.Cell{X: from.X - 1, Y: from.Y, Z: from.Z + 1}
	southWest := client.Cell{X: from.X - 1, Y: from.Y + 1, Z: from.Z}
	southEast := client.Cell{X: from.X, Y: from.Y + 1, Z: from.Z - 1}
	east := client.Cell{X: from.X + 1, Y: from.Y, Z: from.Z - 1}

	if to == northEast {
		return NorthEast
	}
	if to == northWest {
		return NorthWest
	}
	if to == west {
		return West
	}
	if to == southWest {
		return SouthWest
	}
	if to == southEast {
		return SouthEast
	}
	if to == east {
		return East
	}

	panic("smth broken")
}
