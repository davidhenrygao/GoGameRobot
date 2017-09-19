package main

import (
	"fmt"
)

type Changename struct {
	Name string `json:"name"`
}

type Changenameresp struct {
}

type ChangenameResp struct {
	Resp Changenameresp `json:"body"`
}

func ChangenameHandler(b []byte) {
	body := new(EchoResp)
	err := Decode(b, body)
	if err != nil {
		return
	}
	//resp := body.Resp
	fmt.Printf("Receive Changename resp: success.\n")
}

func SendChangename(ch chan<- interface{}) {
	req := Request{Cmd: CHANGENAME, Req: Changename{Name: "DHRobot"}}
	ch <- req
}
