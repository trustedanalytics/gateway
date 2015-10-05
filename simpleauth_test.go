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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validToken   string = "VEVTVA=="
	invalidToken string = "VEVVVA=="
)

func TestDecode(t *testing.T) {
	assert.Equal(t, "TEST", decode(validToken), "Correct token must be decoded correctly")
	assert.Equal(t, "", decode("="), "Incorrect token must be decoded as empty string")
}

func TestSimpleAuth_Validate(t *testing.T) {
	simpleAuth := NewSimpleAuth(validToken)

	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Error(err)
	}

	req.Header.Del("Authorization")
	assert.False(t, simpleAuth.Validate(req), "Requests must have an Authorization header")

	// Test valid request with valid token
	req.Header.Set("Authorization", "Bearer "+validToken)
	assert.True(t, simpleAuth.Validate(req), "Valid request with valid token should be accepted")

	// Test valid request with invalid token
	req.Header.Set("Authorization", "Bearer "+invalidToken)
	assert.False(t, simpleAuth.Validate(req), "Valid request with invalid token should be rejected")

	// Test invalid request with valid token
	req.Header.Set("Authorization", validToken)
	assert.False(t, simpleAuth.Validate(req), "Invalid request with valid token should be rejected")

	// Test invalid request with invalid token
	req.Header.Set("Authorization", invalidToken)
	assert.False(t, simpleAuth.Validate(req), "Invalid request with invalid token should be rejected")
}
