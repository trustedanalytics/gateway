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
	"encoding/json"
	"log"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

// NewMessage factors new message
func NewMessage(body string) *Message {
	return &Message{
		ID:   uuid.New(),
		On:   time.Now().UTC(),
		Body: body,
	}
}

// Message holding type for the message payload
type Message struct {

	// ID is a uuid v4 id of that message
	ID string `json:"id"`

	// On represents when it was received
	On time.Time `json:"on"`

	// Body represents message content
	Body string `json:"body"`
}

// ToBytes converts content of the current message into byte array
func (m *Message) ToBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Printf("unable to marshal: %v", err.Error())
	}
	return b
}
