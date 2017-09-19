package main

import (
	"fmt"
	"time"
)

type Heartbeat struct {
}

type Heartbeatresp struct {
	Servertime uint64 `json:"servertime"`
}

type HeartbeatResp struct {
	Resp Heartbeatresp `json:"body"`
}

func HeartbeatHandler(b []byte) {
	body := new(HeartbeatResp)
	err := Decode(b, body)
	if err != nil {
		return
	}
	resp := body.Resp
	fmt.Printf("Receive Heartbeat resp: severtime(%d).\n", resp.Servertime)
}

func HeartbeatTicker(ch chan<- interface{}, wch chan<- interface{}) {
	ticker := time.NewTicker(2 * time.Second)
	defer func() {
		ticker.Stop()
		wch <- struct{}{}
	}()
	for now := range ticker.C {
		if len(ch) == 0 {
			req := Request{Cmd: HEARTBEAT, Req: Heartbeat{}}
			ch <- req
			fmt.Printf("Generate Heartbeat request(%#v).\n", now)
		}
	}
}
