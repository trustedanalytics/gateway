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
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// ParseInt parses passed string into an Integer
func ParseInt(s string, d int) int {
	if len(s) < 1 {
		return d
	}
	//strconv.Btoi64
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		log.Fatalf("unable to parse int from %s: %v", s, err)
		return d
	}
	return int(v)
}

// GetEnvVarAsString wrapper for env variable with defaults
func GetEnvVarAsString(k, d string) string {
	if len(k) < 1 {
		return d
	}
	s := os.Getenv(k)
	if len(s) < 1 {
		return d
	}
	return s
}

// GetEnvVarAsInt wrapper utility for env variabler as string
func GetEnvVarAsInt(k string, d int) int {
	s := GetEnvVarAsString(k, "")
	if len(s) < 1 {
		return d
	}
	return ParseInt(s, d)
}

// SetWithEnvVar sets variable to the string value of the env variable envVariable
// if it is set. If no env variable by the name envVariable is found, variable
// is not changed.
func SetWithStringEnvVar(envVariable string, variable *string) {
	envVarVal := os.Getenv(envVariable)
	if len(envVarVal) > 0 {
		*variable = envVarVal
	}
}

// GetEnvVarAsBool wrapper for env variable as Bool
func GetEnvVarAsBool(k string, d bool) bool {
	s := GetEnvVarAsString(k, "")
	if len(s) < 1 {
		return d
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatalf("unable to parse bool from %s: %v", s, err)
		return d
	}
	return v
}

// Trace outputs self to object as string
func Trace(str string, o interface{}) {
	if args.Trace {
		objStr, err := ToString(o)
		if err != nil {
			log.Printf("unable to marshal: %v", err.Error())
			return
		}
		log.Printf("%s: %s", str, fmt.Sprintln(string(objStr)))
	}
}

// ToString returns its representaiton as string
func ToString(o interface{}) (string, error) {
	objStr, err := json.Marshal(o)
	if err != nil {
		log.Printf("unable to marshal: %v", o)
		log.Panicln(err)
		return "", err
	}
	return fmt.Sprintln(string(objStr)), nil
}

// GetNowInUtc wrapper for UTC timestamp
func GetNowInUtc() time.Time {
	return time.Now().UTC()
}

// GetTime wrapper for time formater
func GetTime(f string) string {
	if len(f) < 1 {
		f = time.RFC850
	}
	return fmt.Sprintln(GetNowInUtc().Format(f))
}
