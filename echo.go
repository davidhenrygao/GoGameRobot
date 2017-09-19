package main

import (
	"fmt"
)

type Echo struct {
	Msg string `json:"msg"`
}

type Echoresp struct {
	Msg string `json:"msg"`
}

type EchoResp struct {
	Resp Echoresp `json:"body"`
}

func EchoHandler(b []byte) {
	body := new(EchoResp)
	err := Decode(b, body)
	if err != nil {
		return
	}
	resp := body.Resp
	fmt.Printf("Receive Echo resp: echo back(%s).\n", resp.Msg)
}

func SendEcho(ch chan<- interface{}) {
	req := Request{Cmd: ECHO, Req: Echo{Msg: "Hello skynet!"}}
	ch <- req
}
