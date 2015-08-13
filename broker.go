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
	"errors"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

func newBroker() *broker {
	clients := make(map[int64]*handler, 5)
	addCh := make(chan *handler, 5)
	delCh := make(chan *handler)
	doneCh := make(chan bool)
	errCh := make(chan error)
	authV := NewAuth(args.Server.Token)

	return &broker{
		clients,
		addCh,
		delCh,
		doneCh,
		errCh,
		authV,
	}
}

type broker struct {
	clients map[int64]*handler
	addCh   chan *handler
	delCh   chan *handler
	doneCh  chan bool
	errCh   chan error
	authVal *Auth
}

func (s *broker) add(c *handler) { s.addCh <- c }
func (s *broker) del(c *handler) { s.delCh <- c }
func (s *broker) done()          { s.doneCh <- true }
func (s *broker) err(err error)  { s.errCh <- err }
func (s *broker) listen() {

	onConnected := func(ws *websocket.Conn) {

		// make sure closes cleanly
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		// create a new producer client per connection
		if s.authVal.Valid(ws.Request()) {
			handler := newClient(ws, s)
			s.add(handler)
			handler.listen()
		} else {
			log.Println("Invalid token")
			s.errCh <- errors.New("Invalid token")
		}

	}

	onRequest := func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{
			Handler: websocket.Handler(onConnected),
		}
		s.ServeHTTP(w, req)
	}

	http.HandleFunc(args.Server.Root, onRequest)

	for {
		select {
		case c := <-s.addCh:
			s.clients[c.id] = c
			if args.Trace {
				log.Printf("app:%d clients:%d", args.Index, len(s.clients))
			}
		case c := <-s.delCh:
			delete(s.clients, c.id)
			if args.Trace {
				Trace("handler deleted", c)
			}
		case err := <-s.errCh:
			log.Println("error:", err.Error())
		case <-s.doneCh:
			return
		}
	}
}
