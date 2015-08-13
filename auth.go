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
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

// NewAuth factors new authentication object
func NewAuth(token string) *Auth {
	a := &Auth{token: decode(token)}
	a.enabled = len(token) > 0
	return a
}

// Auth is the type representing auth imp
type Auth struct {
	enabled bool
	token   string
}

// Valid validates the authentication from HTTP request
func (a *Auth) Valid(req *http.Request) bool {
	if !a.enabled {
		return true
	}
	auths, _ := req.Header["Authorization"]
	if len(auths) != 1 {
		log.Println("missing Authorization")
		return false
	}
	tokens := strings.Split(auths[0], " ")
	if len(tokens) != 2 || tokens[0] != "Bearer" {
		log.Printf("invalid auth type: %s", tokens)
		return false
	}
	token := decode(tokens[1])
	return a.token == token
}

func decode(val string) string {
	raw, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		log.Printf("unable to decode: %s", val)
		return ""
	}
	return string(raw)
}
