// Copyright 2012-2019 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"time"
)

const clusterID string = "devtron-stan"
const clientID string = "devtron-nats-testing-1001"
const topic string = "foo"

func main() {
	fmt.Println("subscriber main starts ....")
	nc, err := nats.Connect("nats://127.0.0.1:4222", nats.ReconnectWait(10*time.Second), nats.MaxReconnects(100))
	if err != nil {
		fmt.Println(err)
	}
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		fmt.Println(err)
	}
	// Simple Synchronous Publisher
	/*err = sc.Publish("foo", []byte("Hello Viki ..")) // does not return until an ack has been received from NATS Streaming
	if err != nil {
		fmt.Println(err)
	}*/

	fmt.Println(" START Push......")
	for i := 1; i <= 20; i++ {
		msg := fmt.Sprintf("%d noted at: %s", i, time.Now().String())
		err = sc.Publish(topic, []byte(msg))
		if err != nil {
			fmt.Println("err in publishing msg", "err", err)
		}
		fmt.Printf("Sending... %s\n", msg)
	}

	// Simple Async Subscriber
	sub1, err := sc.QueueSubscribe(topic, "foo", func(m *stan.Msg) {
		fmt.Printf("Received a message on sub 1: %s\n", string(m.Data))
		time.Sleep(time.Duration(1) * time.Second)
	}, stan.DurableName("foo"), stan.AckWait(time.Duration(300)*time.Second), stan.StartWithLastReceived(), stan.SetManualAckMode())
	if err != nil {
		fmt.Println(err)
	}

	// Wait for 5 second than REOPEN again
	//time.Sleep(time.Duration(5) * time.Second)
	sub2, err := sc.QueueSubscribe(topic, "foo", func(m *stan.Msg) {
		fmt.Printf("Received a message on sub 2: %s\n", string(m.Data))
		time.Sleep(time.Duration(1) * time.Second)
	}, stan.DurableName("foo"), stan.AckWait(time.Duration(300)*time.Second), stan.StartWithLastReceived(), stan.SetManualAckMode())
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Duration(5) * time.Second)
	// Unsubscribe
	_ = sub1.Unsubscribe()


	time.Sleep(time.Duration(5) * time.Second)
	// Unsubscribe
	_ = sub2.Unsubscribe()

	// Close connection
	_ = sc.Close()
	nc.Close()
	fmt.Println("all done")
}
