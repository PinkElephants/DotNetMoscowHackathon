package client

type Login struct {
	Login    string
	Password string
}

type Start struct {
	Map string
}

type Turn struct {
	Direction    string
	Acceleration int
}
