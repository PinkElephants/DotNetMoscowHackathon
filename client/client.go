package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	host     = "http://51.158.109.80:5000"
	loginURL = host + "/raceapi/Auth/Login"
	raceURL  = host + "/raceapi/race"
	helpURL  = host + "/raceapi/help/math"
	mapName  = "test"
)

type Client struct {
	http    http.Client
	token   Token
	session string
}

func NewClient() *Client {
	return &Client{
		token: Token{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJDcnlwdG9FbGVwaGFudHMiLCJqdGkiOiIyNTNmZTgxOC02MTM5LTQ1OTctOTY5Ni02ZjBhM2Y4ODhkNjYiLCJodHRwOi8vc2NoZW1hcy54bWxzb2FwLm9yZy93cy8yMDA1LzA1L2lkZW50aXR5L2NsYWltcy9uYW1laWRlbnRpZmllciI6IjYzNTRiMzgzLWM5YWMtNGRhNC05MWQ2LTA1MGE4YTMwNjI1YSIsInJvbCI6ImFwaV9hY2Nlc3MiLCJleHAiOjE1NTcxMzQ0NjIsImlzcyI6Im1za2RvdG5ldCIsImF1ZCI6Im1za2RvdG5ldCJ9.GJCqbvkDlSDig-2ezFoqB95pZJcxnmlqQCaqTDDTf3c"},
	}
}

func (c *Client) Login() {
	if len(c.token.Token) != 0 {
		return
	}
	login, err := json.Marshal(Login{
		Login:    "CryptoElephants",
		Password: "46ZHlr",
	})

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer([]byte(login)))
	checkErr(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	checkErr(err)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &c.token)
	checkErr(err)
}

func (c *Client) Help() Help {
	req, err := http.NewRequest("GET", helpURL, nil)
	checkErr(err)
	resp, err := c.http.Do(req)
	checkErr(err)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		panic(string(body))
	}

	help := Help{}
	err = json.Unmarshal(body, &help)

	checkErr(err)
	return help
}

func (c *Client) Start() ServerInfo {
	start, err := json.Marshal(Start{
		Map: mapName,
	})

	req := c.request("POST", raceURL, start)
	resp, err := c.http.Do(req)
	checkErr(err)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		panic(string(body))
	}

	info := ServerInfo{}
	err = json.Unmarshal(body, &info)
	checkErr(err)
	c.session = info.SessionID
	return info
}

func (c *Client) UpdateUI() {
	req := c.request("GET", raceURL+"?="+c.session, nil)
	resp, err := c.http.Do(req)
	checkErr(err)
	ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
}

func (c *Client) Turn(t Turn) TurnResult {
	turn, err := json.Marshal(t)

	req := c.request("PUT", c.turnUrl(), turn)
	resp, err := c.http.Do(req)
	checkErr(err)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		panic(string(body))
	}

	res := TurnResult{}
	err = json.Unmarshal(body, &res)
	checkErr(err)
	return res
}

func (c *Client) turnUrl() string {
	return raceURL + "/" + c.session
}

func (c *Client) request(method, url string, body []byte) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	checkErr(err)
	req.Header.Set("Content-Type", "application/json")
	if len(c.token.Token) != 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token.Token))
	}

	return req
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
