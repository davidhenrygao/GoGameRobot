package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

//ConnContext Client connection context.
type ConnContext struct {
	c  net.Conn
	ch chan interface{}
}

//SetConn connection setter
func (rc *ConnContext) SetConn(conn net.Conn) {
	rc.c = conn
}

//SetChan ConnContext channel setter
func (rc *ConnContext) SetChan(ch chan interface{}) {
	rc.ch = ch
}

func (rc *ConnContext) GetChan() chan interface{} {
	return rc.ch
}

//SendPackage package sender
func (rc *ConnContext) SendPackage(b []byte) error {
	pkg, err := GenPackage(b)
	if err != nil {
		fmt.Printf("GenPackage error: %s.\n", err)
		return err
	}
	_, err = rc.c.Write(pkg)
	if err != nil {
		fmt.Printf("Connection write error: %s.\n", err)
		return err
	}
	return nil
}

//RecvPackage receive package bytes
func (rc *ConnContext) RecvPackage() ([]byte, error) {
	head := make([]byte, 2)
	_, err := io.ReadFull(rc.c, head)
	if err != nil {
		fmt.Printf("RecvPackage read head error: %s.\n", err)
		return nil, err
	}
	ul := binary.BigEndian.Uint16(head)
	body := make([]byte, ul)
	_, err = io.ReadFull(rc.c, body)
	if err != nil {
		fmt.Printf("RecvPackage read msg body error: %s.\n", err)
		return nil, err
	}
	return body, nil
}

//GenPackage pack a header to json string
func GenPackage(b []byte) ([]byte, error) {
	l := len(b)
	if l >= 2<<8 {
		err := fmt.Errorf("length(%d) exceeds %d", l, 2<<8)
		return nil, err
	}
	ul := uint16(l)
	rb := make([]byte, 2+ul)
	binary.BigEndian.PutUint16(rb[:2], ul)
	copy(rb[2:], b)
	return rb, nil
}

//Encode to Json
func Encode(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("Json marshal error: %s.\n", err)
		return nil, err
	}
	return b, nil
}

//Decode form Json, v must be a pointer of the response message struct
func Decode(b []byte, v interface{}) error {
	err := json.Unmarshal(b, v)
	if err != nil {
		fmt.Printf("Json unmarshal error: %s.\n", err)
		return err
	}
	return nil
}

func ContextStart(context *ConnContext) {
	var b []byte
	var err error
	ch := context.GetChan()
	for req := range ch {
		r, ok := req.(Request)
		if !ok {
			fmt.Printf("Type assertion failed: req(%#v) not Request struct!", req)
			continue
		}
		b, err = Encode(req)
		if err != nil {
			continue
		}
		err = context.SendPackage(b)
		if err != nil {
			continue
		}
		b, err = context.RecvPackage()
		if err != nil {
			continue
		}
		HandleCmdResp(r.Cmd, b)
	}
}

func main() {
	conn, err := net.Dial("tcp", "192.168.2.188:10086")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	context := new(ConnContext)
	context.SetConn(conn)
	ch := make(chan interface{}, 100)
	context.SetChan(ch)

	wch := make(chan interface{})
	go HeartbeatTicker(ch, wch)
	go ContextStart(context)
	WorkSequence(ch)

	<-wch
}
