package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"time"
)

const clusterID string = "devtron-stan"
const clientID string = "devtron-nats-n-testing-1001"
const subject string = "nishant-test-2"

func main() {
	fmt.Println("subscriber 1 invoked ..")
	nc, err := nats.Connect("nats://127.0.0.1:4222", nats.ReconnectWait(10*time.Second), nats.MaxReconnects(100))
	if err != nil {
		fmt.Print(err)
	}
	defer nc.Close()
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		fmt.Println(err)
	}
	_, _ = sc.QueueSubscribe(subject, "ge", func(msg *stan.Msg) {
		fmt.Printf("Received a message on sub 1: %s\n", string(msg.Data))
		time.Sleep(1 * time.Second)
		err := msg.Ack()
		if err != nil {
			fmt.Println("ack err" + err.Error())
		}
	}, stan.StartWithLastReceived(),
		stan.AckWait(300*time.Second),
		stan.SetManualAckMode(),
		stan.DurableName("my-durable"))

	// Simple Synchronous Publisher
	time.Sleep(time.Duration(50) * time.Minute)
	_ = sc.Close()
}
