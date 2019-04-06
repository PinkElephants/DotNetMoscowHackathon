package client

type Token struct {
	Token string
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
