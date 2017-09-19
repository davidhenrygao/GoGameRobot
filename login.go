package main

import (
	"fmt"
)

type Login struct {
	Account string `json:"account"`
	Passwd  string `json:"passwd"`
}

type Loginresp struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type LoginResp struct {
	Resp Loginresp `json:"body"`
}

func LoginHandler(b []byte) {
	body := new(LoginResp)
	err := Decode(b, body)
	if err != nil {
		return
	}
	resp := body.Resp
	fmt.Printf("Receive Login resp: playerinfo(%#v).\n", resp)
}

func SendLogin(ch chan<- interface{}) {
	req := Request{Cmd: LOGIN, Req: Login{Account: "TestRobot", Passwd: "12345"}}
	ch <- req
}
