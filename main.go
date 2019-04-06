package main

import (
	"github.com/PinkElephants/DotNetMoscowHackathon/client"
)

func main() {
	c := client.NewClient()
	c.Login()
	c.Start()

}
