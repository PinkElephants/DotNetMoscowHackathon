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
	b.Start(c.Start())

	for {
		res := c.Turn(b.Turn())
		//c.UpdateUI()
		if res.Status == calc.Hungry || res.Status == calc.Punished || res.Status == calc.HappyAsInsane {
			return
		}
		b.Result(res)
	}
}
