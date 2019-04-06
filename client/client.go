package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	host     = "http://51.15.100.12:5000"
	loginURL = host + "/raceapi/Auth/Login"
	raceURL  = host + "/raceapi/race"
	mapName  = "test"
)

type Client struct {
	http  http.Client
	token Token
	info  ServerInfo
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

func (c *Client) Start() {
	start, err := json.Marshal(Start{
		Map: mapName,
	})

	req := c.request(raceURL, start)
	resp, err := c.http.Do(req)
	checkErr(err)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &c.info)
	checkErr(err)
}

func (c *Client) Turn() {
	req := c.request(raceURL, []byte("test"))
	resp, err := c.http.Do(req)
	checkErr(err)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Turn: ", string(body))
}

func (c *Client) request(url string, body []byte) *http.Request {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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
