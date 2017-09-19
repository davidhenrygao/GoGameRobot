package main

import (
	"fmt"
)

const (
	HEARTBEAT = 1
	ECHO      = 2

	LOGIN      = 100
	CHANGENAME = 101
)

const (
	SUCCESS                      = 0
	INTERNAL                     = 1
	UNKNOWN_CMD                  = 2
	PROTO_UNSERIALIZATION_FAILED = 3

	//1000-9999 are used for login
	ACCOUNT_ALREADY_EXIST = 1000
	ACCOUNT_NOT_EXIST     = 1001

	//10000-19999 are used for agent
	PLAYER_ID_NOT_EXIT = 10000
)

type CODE struct {
	Code uint `json:"code"`
}

type Request struct {
	Cmd uint        `json:"code"`
	Req interface{} `json:"body"`
}

func WorkSequence(ch chan<- interface{}) {
	SendLogin(ch)
	SendChangename(ch)
	SendEcho(ch)
	SendLogin(ch)
}

func HandleCmdResp(cmd uint, b []byte) {
	code := new(CODE)
	err := Decode(b, code)
	if err != nil {
		return
	}
	if code.Code != SUCCESS {
		fmt.Printf("Request(%d) response failed: error code(%d).\n", cmd, code.Code)
		return
	}
	switch cmd {
	case HEARTBEAT:
		HeartbeatHandler(b)
	case ECHO:
		EchoHandler(b)
	case LOGIN:
		LoginHandler(b)
	case CHANGENAME:
		ChangenameHandler(b)
	default:
		fmt.Printf("Unknown command %d.\n", cmd)
	}
}
