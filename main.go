package main

import (
	"github.com/PinkElephants/DotNetMoscowHackathon/bot"
	"github.com/PinkElephants/DotNetMoscowHackathon/client"
)

func main() {
	b := calc.NewBot()
	c := client.NewClient()
	c.Login()
	b.Help = c.Help()
	c.Start()
}
