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
	"github.com/Shopify/sarama"
	"log"
	"time"
)

var (
	qProducer sarama.AsyncProducer
)

func queueInit() {

	config := sarama.NewConfig()

	config.ClientID = args.ID

	// Acks
	if args.Pub.Ack {
		config.Producer.RequiredAcks = sarama.WaitForAll
	} else {
		config.Producer.RequiredAcks = sarama.WaitForLocal
	}

	// Compress
	if args.Pub.Compress {
		config.Producer.Compression = sarama.CompressionSnappy
	} else {
		config.Producer.Compression = sarama.CompressionNone
	}

	// Flush Intervals
	if args.Pub.FlushFreq > 0 {
		config.Producer.Flush.Frequency = time.Duration(args.Pub.FlushFreq) * time.Second
	} else {
		config.Producer.Flush.Frequency = 1 * time.Second
	}

	producer, err := sarama.NewAsyncProducer(args.Pub.URI, config)
	if err != nil {
		log.Fatalln("Failed to start Kafka producer:", err)
	}

	qProducer = producer

}

// Start fires a publisher listener
func produce(in <-chan *Message) {

	for {
		msg := <-in
		select {
		case qProducer.Input() <- &sarama.ProducerMessage{
			Topic: args.Pub.Topic,
			Key:   nil,
			Value: sarama.StringEncoder(msg.ToBytes()),
		}:
			if args.Trace {
				log.Printf("Queue[%s] < %s", args.Pub.Topic, msg)
			}
		case err := <-qProducer.Errors():
			log.Printf("Error on queue send for [%s]: %v", args.Pub.Topic, err.Err)
		}
	}

}
