package client

import "math"

type Token struct {
	Token string
}

type Help struct {
	MaxSpeed        int `json:"MaxSpeed"`
	MinSpeed        int `json:"MinSpeed"`
	MaxAcceleration int `json:"MaxAcceleration"`
	DriftsAngles    []struct {
		Angle          int `json:"Angle"`
		MaxSpeed       int `json:"MaxSpeed"`
		SpeedDownShift int `json:"SpeedDownShift"`
	} `json:"DriftsAngles"`
	MinCanyonSpeed          int `json:"MinCanyonSpeed"`
	MaxDuneSpeed            int `json:"MaxDuneSpeed"`
	BaseTurnFuelWaste       int `json:"BaseTurnFuelWaste"`
	DriftFuelMultiplier     int `json:"DriftFuelMultiplier"`
	FullSpeedFuelMultiplier int `json:"FullSpeedFuelMultiplier"`
	Angles                  []struct {
		Direction string `json:"Direction"`
		Angle     int    `json:"Angle"`
	} `json:"Angles"`
	LocationDeltas []struct {
		Direction string `json:"Direction"`
		Delta     struct {
			Dx int `json:"Dx"`
			Dy int `json:"Dy"`
			Dz int `json:"Dz"`
		} `json:"Delta"`
	} `json:"LocationDeltas"`
}

type Car struct {
	X, Y, Z      int
	Acceleration int
	Direction    string
	Heading      string
	Speed        int
	Fuel         int
}

type Cell struct {
	X, Y, Z      int
	Type         string
	DistToCar    int
	DistToTarget int
	Visible      bool
}

func (c *Cell) DistanceFrom(origin Cell) int {
	res := (math.Abs(float64(c.X)-float64(origin.X)) + math.Abs(float64(c.Y)-float64(origin.Y)) + math.Abs(float64(c.Z)-float64(origin.Z))) / 2
	return int(res)
}

type ServerInfo struct {
	SessionID        string `json:"SessionId"`
	PlayerID         string `json:"PlayerId"`
	CurrentDirection string `json:"CurrentDirection"`
	CurrentLocation  struct {
		X int `json:"X"`
		Y int `json:"Y"`
		Z int `json:"Z"`
	} `json:"CurrentLocation"`
	Finish struct {
		X int `json:"X"`
		Y int `json:"Y"`
		Z int `json:"Z"`
	} `json:"Finish"`
	Radius         int    `json:"Radius"`
	CurrentSpeed   int    `json:"CurrentSpeed"`
	PlayerStatus   string `json:"PlayerStatus"`
	NeighbourCells []struct {
		Item1 struct {
			X int `json:"X"`
			Y int `json:"Y"`
			Z int `json:"Z"`
		} `json:"Item1"`
		Item2 string `json:"Item2"`
	} `json:"NeighbourCells"`
	Fuel int `json:"Fuel"`
}

func (s *ServerInfo) Cells() []Cell {
	cells := make([]Cell, len(s.NeighbourCells))
	for i, c := range s.NeighbourCells {
		cells[i] = Cell{
			X:    c.Item1.X,
			Y:    c.Item1.Y,
			Z:    c.Item1.Z,
			Type: c.Item2,
		}
	}
	return cells
}

func (s *ServerInfo) Car() Car {
	return Car{
		X:         s.CurrentLocation.X,
		Y:         s.CurrentLocation.Y,
		Z:         s.CurrentLocation.Z,
		Direction: s.CurrentDirection,
		Fuel:      s.Fuel,
	}
}

type TurnResult struct {
	Command struct {
		Location struct {
			X int `json:"X"`
			Y int `json:"Y"`
			Z int `json:"Z"`
		} `json:"Location"`
		Acceleration      int    `json:"Acceleration"`
		MovementDirection string `json:"MovementDirection"`
		Heading           string `json:"Heading"`
		Speed             int    `json:"Speed"`
		Fuel              int    `json:"Fuel"`
	} `json:"Command"`
	VisibleCells []struct {
		Item1 struct {
			X int `json:"X"`
			Y int `json:"Y"`
			Z int `json:"Z"`
		} `json:"Item1"`
		Item2 string `json:"Item2"`
	} `json:"VisibleCells"`
	Location struct {
		X int `json:"X"`
		Y int `json:"Y"`
		Z int `json:"Z"`
	} `json:"Location"`
	ShortestWayLength int    `json:"ShortestWayLength"`
	Speed             int    `json:"Speed"`
	Status            string `json:"Status"`
	Heading           string `json:"Heading"`
	FuelWaste         int    `json:"FuelWaste"`
}

func (s *TurnResult) Cells() []Cell {
	cells := make([]Cell, len(s.VisibleCells))
	for i, c := range s.VisibleCells {
		cells[i] = Cell{
			X:    c.Item1.X,
			Y:    c.Item1.Y,
			Z:    c.Item1.Z,
			Type: c.Item2,
		}
	}
	return nil
}

func (s *TurnResult) Car() Car {
	return Car{
		X:            s.Location.X,
		Y:            s.Location.Y,
		Z:            s.Location.Z,
		Direction:    s.Command.MovementDirection,
		Acceleration: s.Command.Acceleration,
		Heading:      s.Command.Heading,
		Speed:        s.Command.Speed,
	}
}
