package calc

import (
	"fmt"
	"math/rand"
	"sort"

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

var angles = []string{
	East,
	West,
	NorthEast,
	NorthWest,
	SouthEast,
	SouthWest,
}

const (
	Empty         = "Empty"
	Rock          = "Rock"
	DangerousArea = "DangerousArea"
	Pit           = "Pit"
)

type Bot struct {
	Help client.Help

	info   client.ServerInfo
	result client.TurnResult

	turn       client.Turn
	car        client.Car
	cellsIndex [][][]client.Cell

	wasAcceleratedPrevious bool
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) Start(info client.ServerInfo) {
	b.info = info
	b.car = info.Car()
	b.allocateCells()
	b.updateCells(info.Cells())
	b.scan()

	b.makeTurn()
}

func (b *Bot) Result(result client.TurnResult) {
	b.car = result.Car()
	b.updateCells(result.Cells())
	b.scan()

	b.makeTurn()
}

func (b *Bot) makeTurn() {
	path := b.happyPath()

	goTo := b.selectGoTo(path)
	b.turn = client.Turn{
		Direction:    goTo,
		Acceleration: b.acceleration(b.cellFromAndAngle(b.carCell(), goTo)),
	}
}

func (b *Bot) Turn() client.Turn {
	return b.turn
}

func (b *Bot) happyPath() []client.Cell {
	var path []client.Cell

	current := b.closestToTarget()
	for {
		best := current
		b.iterNeighbors(current, func(c client.Cell) {
			if c.Type == Rock || !c.Visible {
				return
			}
			if c.DistToCar < best.DistToCar {
				best = c
			}
		})
		if best.DistToCar == 0 {
			break
		}
		current = best
		path = append(path, current)
	}
	return path
}

func (b *Bot) allocateCells() {
	b.cellsIndex = make([][][]client.Cell, b.info.Radius*2)
	for i := range b.cellsIndex {
		b.cellsIndex[i] = make([][]client.Cell, b.info.Radius*2)
		for j := range b.cellsIndex[i] {
			b.cellsIndex[i][j] = make([]client.Cell, b.info.Radius*2)
		}
	}
}

func (b *Bot) updateCells(cells []client.Cell) {
	carCell := b.cell(b.car.X, b.car.Y, b.car.Z)
	carCell.Visible = true
	carCell.Visited = true
	b.setCell(b.car.X, b.car.Y, b.car.Z, carCell)
	for _, c := range cells {
		c.Visible = true
		b.setCell(c.X, c.Y, c.Z, c)
	}
}

func (b *Bot) acceleration(goTo client.Cell) int {
	safeSpeed := ((b.Help.MinCanyonSpeed - b.Help.MaxDuneSpeed) / 2) + b.Help.MaxDuneSpeed

	if goTo.Type == Pit {
		return b.Help.MinCanyonSpeed - b.car.Speed
	}
	if goTo.Type == DangerousArea {
		return b.car.Speed - b.Help.MaxDuneSpeed
	}
	//if goTo.Type == DangerousArea {
	//	acceleration := b.car.Speed - b.Help.MaxDuneSpeed
	//	b.car.Speed += acceleration
	//	return acceleration
	//}

	// if goTo.Type != Pit && b.wasAcceleratedPrevious {
	// 	b.wasAcceleratedPrevious = false
	// }
	//
	// if goTo.Type == Pit {
	// 	b.wasAcceleratedPrevious = true
	// 	acceleration := b.Help.MinCanyonSpeed - b.car.Speed
	// 	b.car.Speed += acceleration
	// 	return acceleration
	// }
	// if goTo.Type == DangerousArea {
	// 	acceleration := b.car.Speed - b.Help.MaxDuneSpeed
	// 	b.car.Speed += acceleration
	// 	return acceleration
	// }

	return safeSpeed - b.car.Speed
}

func (b *Bot) scan() {
	b.iterAll(func(c client.Cell) {
		if c.Type == Rock {
			return
		}
		if !c.Visible {
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

		c.NearWall = b.nearWall(c)

		b.setCell(c.X, c.Y, c.Z, c)
	})
}

func (b *Bot) closestToTarget() client.Cell {
	dist := b.cell(b.car.X, b.car.Y, b.car.Z)

	b.iterAll(func(c client.Cell) {
		if c.Type == Rock || !c.Visible {
			return
		}
		if c.DistToTarget < dist.DistToTarget {
			dist = c
		}
	})

	return dist
}

func (b *Bot) iterNeighbors(cell client.Cell, f func(c client.Cell)) {
	f(b.cell(cell.X-1, cell.Y+1, cell.Z))
	f(b.cell(cell.X, cell.Y+1, cell.Z-1))
	f(b.cell(cell.X+1, cell.Y, cell.Z-1))
	f(b.cell(cell.X+1, cell.Y-1, cell.Z))
	f(b.cell(cell.X, cell.Y-1, cell.Z+1))
	f(b.cell(cell.X-1, cell.Y, cell.Z+1))
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

func (b *Bot) cell(x, y, z int) client.Cell {
	r := b.info.Radius
	if x > r || x < -r {
		panic("ups")
	}
	if y > r || y < -r {
		panic("ups")
	}
	if z > r || z < -r {
		panic("ups")
	}
	return b.cellsIndex[x+r][y+r][z+r]
}

func (b *Bot) setCell(x, y, z int, c client.Cell) {
	r := b.info.Radius
	b.cellsIndex[x+r][y+r][z+r] = c
}

func angle(from client.Cell, to client.Cell) string {
	northEast := client.Cell{X: from.X + 1, Y: from.Y, Z: from.Z - 1}
	northWest := client.Cell{X: from.X, Y: from.Y + 1, Z: from.Z - 1}
	west := client.Cell{X: from.X - 1, Y: from.Y + 1, Z: from.Z}
	southWest := client.Cell{X: from.X - 1, Y: from.Y, Z: from.Z + 1}
	southEast := client.Cell{X: from.X, Y: from.Y - 1, Z: from.Z + 1}
	east := client.Cell{X: from.X + 1, Y: from.Y - 1, Z: from.Z}

	if to.Equal(northEast) {
		return NorthEast
	}
	if to.Equal(northWest) {
		return NorthWest
	}
	if to.Equal(west) {
		return West
	}
	if to.Equal(southWest) {
		return SouthWest
	}
	if to.Equal(southEast) {
		return SouthEast
	}
	if to.Equal(east) {
		return East
	}

	fmt.Println("wrong goto")
	return angles[rand.Int()%6]
}

func (b *Bot) nearWall(cell client.Cell) bool {
	var near bool
	b.iterNeighbors(cell, func(c client.Cell) {
		if c.Type == Rock {
			near = true
		}
	})
	return near
}

func (b *Bot) selectGoTo(path []client.Cell) string {
	type candidate struct {
		angle    string
		priority int
	}

	var heading string
	if len(b.result.Heading) != 0 {
		heading = b.result.Heading
	} else {
		if len(path) == 0 {
			heading = angles[rand.Int()%6]
		} else {
			happy := path[len(path)-1]
			heading = angle(
				b.carCell(),
				happy,
			)
		}
	}

	cell := b.cellFromAndAngle(b.carCell(), heading)
	if cell.Type != Rock && cell.Visible {
		return heading
	}

	left := heading
	right := heading
	var cands []candidate
	for i := 0; i < 3; i++ {
		left = calcLeft(left)
		right = calcRight(right)
		leftCell := b.cellFromAndAngle(b.carCell(), left)
		rightCell := b.cellFromAndAngle(b.carCell(), right)

		if leftCell.Visible && leftCell.Type != Rock {
			leftPri := 0
			if leftCell.Visited {
				leftPri += 2
			}
			if leftCell.NearWall {
				leftPri += 1
			}
			cands = append(cands, candidate{angle: left, priority: leftPri + i*3})
		}

		if rightCell.Visible && rightCell.Type != Rock {
			rightPri := 0
			if rightCell.Visited {
				rightPri += 2
			}
			if rightCell.NearWall {
				rightPri += 2
			}
			cands = append(cands, candidate{angle: right, priority: rightPri + i*3})
		}
	}

	sort.Slice(cands, func(i, j int) bool {
		return cands[i].priority < cands[j].priority
	})

	if len(cands) == 0 {
		return angles[rand.Int()%6]
	}

	return cands[0].angle
}

func (b *Bot) carCell() client.Cell {
	return b.cell(b.car.X, b.car.Y, b.car.Z)
}

func (b *Bot) setCarCell(cell client.Cell) {
	b.setCell(b.car.X, b.car.Y, b.car.Z, cell)
}

func calcRight(from string) string {
	switch from {
	case East:
		return SouthEast
	case SouthEast:
		return SouthWest
	case SouthWest:
		return West
	case West:
		return NorthWest
	case NorthWest:
		return NorthEast
	case NorthEast:
		return East
	}

	return angles[rand.Int()%6]
}

func calcLeft(from string) string {
	switch from {
	case East:
		return NorthEast
	case NorthEast:
		return NorthWest
	case NorthWest:
		return West
	case West:
		return SouthWest
	case SouthWest:
		return SouthEast
	case SouthEast:
		return East
	}

	return angles[rand.Int()%6]
}

func (b *Bot) cellFromAndAngle(from client.Cell, heading string) client.Cell {
	if heading == NorthEast {
		return b.cell(from.X+1, from.Y, from.Z-1)
	}
	if heading == NorthWest {
		return b.cell(from.X, from.Y+1, from.Z-1)
	}
	if heading == West {
		return b.cell(from.X-1, from.Y+1, from.Z)
	}
	if heading == SouthWest {
		return b.cell(from.X-1, from.Y, from.Z+1)
	}
	if heading == SouthEast {
		return b.cell(from.X, from.Y-1, from.Z+1)
	}
	if heading == East {
		return b.cell(from.X+1, from.Y-1, from.Z)
	}

	panic("smth broken")
}
