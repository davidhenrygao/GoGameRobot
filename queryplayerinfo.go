package main

import (
	"fmt"
)

type Queryplayerinfo struct {
	Id uint64 `json:"id"`
}

type Queryplayerinforesp struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Gold  uint32 `json:"gold"`
	Level uint32 `json:"level"`
	Exp   uint32 `json:"exp"`
}

type QueryplayerinfoResp struct {
	Resp Queryplayerinforesp `json:"body"`
}

func QueryplayerinfoHandler(b []byte) {
	body := new(QueryplayerinfoResp)
	err := Decode(b, body)
	if err != nil {
		return
	}
	resp := body.Resp
	fmt.Printf("Receive query plyaer info resp: playerinfo(%#v).\n", resp)
}

func SendQueryplayerinfo(ch chan<- interface{}) {
	req := Request{Cmd: QUERYPLAYERINFO, Req: Queryplayerinfo{Id: 1}}
	ch <- req
}
