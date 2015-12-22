/**
 * Copyright (c) 2015 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"fmt"
	"io"
	"log"
	"sync/atomic"

	"code.google.com/p/go.net/websocket"
)

const (
	channelBufSize = 100
)

var (
	maxClientID int64
	maxMsgID    int64
	msg         = websocket.Message
)

type handler struct {
	id     int64
	ws     *websocket.Conn
	server *broker
	ch     chan *interface{}
	sender chan *Message
}

// Publisher defines the implementation for publisher
type Publisher interface {
	Config(clientID string, args *PubConfig)
	Start(in <-chan *Message)
}

func newClient(ws *websocket.Conn, s *broker) *handler {
	if ws == nil {
		panic("ws cannot be nil")
	}
	if s == nil {
		panic("server cannot be nil")
	}
	atomic.AddInt64(&maxClientID, 1)
	ch := make(chan *interface{}, channelBufSize)

	h := &handler{
		id:     maxClientID,
		ws:     ws,
		server: s,
		ch:     ch,
		sender: make(chan *Message, 1),
	}

	go produce(h.sender)
	return h
}

func (c *handler) write(msg *interface{}) {
	select {
	case c.ch <- msg:
	default:
		c.server.del(c)
		err := fmt.Errorf("handler %d is disconnected on %s",
			c.id, args.Index)
		c.server.err(err)
	}
}

func (c *handler) conn() *websocket.Conn { return c.ws }
func (c *handler) listen()               { c.listenRead() }
func (c *handler) listenRead() {
	for {
		var m string
		err := msg.Receive(c.ws, &m)
		if err == io.EOF {
			c.server.del(c)
			return
		} else if err != nil {
			c.server.err(err)
		} else {
			if args.Trace {
				atomic.AddInt64(&maxMsgID, 1)
				log.Printf("handler[%d] queued > msg[%d]:%s",
					c.id, maxMsgID, m)
			}
			c.sender <- NewMessage(m)
		}
	}
}
