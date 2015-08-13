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
	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
)

var args = Config{}

func init() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	// load from local file
	loadConfig("./defaults.json", &args)

	args.Trace = GetEnvVarAsBool("GATEWAY_TRACE", false)
	args.Pub.Ack = GetEnvVarAsBool("GATEWAY_ACKS", false)
	args.Pub.Compress = GetEnvVarAsBool("GATEWAY_COMPRESS", true)
	args.Server.Token = os.Getenv("GATEWAY_TOKEN")
	args.Pub.Topic = os.Getenv("GATEWAY_TOPIC")

	var kafkaNodes string = os.Getenv("GATEWAY_QUEUE")

	cf, _ := cfenv.Current()

	if cf != nil {

		Trace("CF", cf)

		args.ID = fmt.Sprintf("%s-%d", cf.ID, cf.Index)
		args.Index = cf.Index
		args.Server.Port = cf.Port
		args.Server.Host = cf.Host
		args.Pub.Topic = cf.Name

		kafka, _ := cf.Services.WithTag("kafka")
		if len(kafka) > 0 {
			kafkaNodes = kafka[0].Credentials["uri"]
		}

	} else {
		log.Println("No CF")
	}

	if len(kafkaNodes) > 0 {
		args.Pub.URI = strings.Split(kafkaNodes, ",")
	}

	Trace("config", args)
}

// ServerConfig represents the Web server configuration holder
type ServerConfig struct {
	Root  string `json:"root,omitempty"`
	Host  string `json:"host,omitempty"`
	Port  int    `json:"port,omitempty"`
	Token string `json:"token,omitempty"`
}

// PubConfig represents the publisher configuration holder
type PubConfig struct {
	URI       []string `json:"uri,omitempty"`
	Topic     string   `json:"topic,omitempty"`
	Ack       bool     `json:"args,acks"`
	Compress  bool     `json:"args,compress"`
	FlushFreq int      `json:"args,flushevery"`
}

// Config represents the root object configuraiton holder
type Config struct {
	ID     string       `json:"id,omitempty"`
	Index  int          `json:"index,omitempty"`
	Trace  bool         `json:"trace,omitempty"`
	Server ServerConfig `json:"server,omitempty"`
	Pub    PubConfig    `json:"publisher,omitempty"`
}

func loadConfig(path string, c *Config) {
	log.Printf("loading config from file: %s ...", path)
	f, err := os.Open(path)
	if err != nil {
		log.Panicf("error while reading config file: %s - %v", path, err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		log.Panicf("error while parsing config file: %s - %v", path, err)
	}
}
