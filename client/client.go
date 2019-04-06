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
)

type Client struct {
	http  http.Client
	token Token
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Login() {
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
	fmt.Println("Token: ", c.token.Token)
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
