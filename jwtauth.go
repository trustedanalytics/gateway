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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtAuth is the type representing JWT auth imp
type JwtAuth struct {
}

type DeviceKeyResponseBody struct {
	PublicKey string `json:"public_key"`
}

// NewJwtAuth creates a new JWT authentication object
func NewJwtAuth() Authenticator {
	return &JwtAuth{}
}

// Validate validates the authentication from HTTP request
func (a *JwtAuth) Validate(req *http.Request) bool {
	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		var alg string = token.Header[jwtAlgFieldName].(string)

		_, methodIsEcdsa := token.Method.(*jwt.SigningMethodECDSA)
		_, methodIsRsa := token.Method.(*jwt.SigningMethodRSA)
		_, methodIsRsaPss := token.Method.(*jwt.SigningMethodRSAPSS)
		if !methodIsEcdsa && !methodIsRsa && !methodIsRsaPss {
			return nil, fmt.Errorf("unexpected signing method: %v", alg)
		}

		iat := token.Claims[iatJWTPayloadFieldName]
		if iat == nil {
			return nil, fmt.Errorf("auth JWT payload must contain an `iat` field")
		}
		issuedAt := time.Unix(int64(iat.(float64)), 0)

		if !isJWTIATAcceptable(issuedAt) {
			return nil, fmt.Errorf("JWT iat not acceptable")
		}

		deviceID := token.Claims[deviceIDJWTPayloadFieldName]
		if deviceID == nil {
			return nil, fmt.Errorf("auth JWT payload must contain a `device_id` field")
		}

		var verifyBytes []byte

		verifyBytes, err := getPublicKeyFromDeviceKeysAPI(deviceID.(string), alg)
		if err != nil {
			return nil, fmt.Errorf("unable to get public key from device keys API: %v", err)
		}

		var verifyKey interface{}

		algIsEs, err := regexp.MatchString("^ES[[:digit:]]+$", alg)
		if err != nil {
			return nil, err
		}

		if algIsEs {
			verifyKey, err = jwt.ParseECPublicKeyFromPEM(verifyBytes)
		} else {
			verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		}
		if err != nil {
			return nil, err
		}

		return verifyKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				log.Println("Malformed token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				log.Println("Token is either expired or not active yet")
			} else {
				log.Println("Couldn't handle this token:", err)
			}
		} else {
			log.Println("Couldn't handle this token:", err)
		}

		return false
	}

	return token.Valid
}

func getPublicKeyFromDeviceKeysAPI(deviceID string, alg string) ([]byte, error) {
	requestURL, err := buildDeviceKeyRequestURL(args.Server.DeviceKeysURI, deviceID, alg)
	if err != nil {
		return nil, fmt.Errorf("unable to build a device key request URL: %v", err)
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("unable to access API to retrieve a public key: %v", err)
	}
	defer resp.Body.Close()

	keyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}

	var respBody DeviceKeyResponseBody
	err = json.Unmarshal(keyBytes, &respBody)
	if err != nil {
		return nil, fmt.Errorf("unable to parse response body: %v", err)
	}
	if len(respBody.PublicKey) == 0 {
		return nil, fmt.Errorf("response body had no public key field")
	}

	keyBytes, err = hex.DecodeString(respBody.PublicKey)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to decode hex string representing public key from response body: %v", err)
	}

	return keyBytes, nil
}

const (
	deviceIDJWTPayloadFieldName string        = "device_id"
	iatJWTPayloadFieldName      string        = "iat"
	jwtAlgFieldName             string        = "alg"
	deviceKeyRequestTimeout     time.Duration = 10 * time.Second
)

var client *http.Client = &http.Client{Timeout: deviceKeyRequestTimeout}

func buildDeviceKeyRequestURL(rawURL string, deviceID string, alg string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("unable to parse the device keys URI: %v", err)
	}
	q := u.Query()
	q.Set(jwtAlgFieldName, alg)
	u.RawQuery = q.Encode()
	u.Path = strings.Replace(u.Path, ":"+deviceIDJWTPayloadFieldName, deviceID, 1)

	return u.String(), nil
}

func isJWTIATAcceptable(issuedAt time.Time) bool {
	return time.Now().Before(
		issuedAt.Add(time.Duration(args.Server.TolerableJWTAge) * time.Minute))
}
