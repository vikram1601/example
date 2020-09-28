package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"time"
)

const clusterID string = "devtron-stan"
const clientID string = "devtron-nats-n-testing-1002"
const subject string = "nishant-test-2"

func main() {
	fmt.Println("producer invoked")
	nc, err := nats.Connect("nats://127.0.0.1:4222", nats.ReconnectWait(10*time.Second), nats.MaxReconnects(100))
	if err != nil {
		fmt.Print(err)
	}
	defer nc.Close()
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		fmt.Println(err)
	}
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				err = sc.Publish(subject, []byte("Hello World "+t.String())) // does not return until an ack has been received from NATS Streaming
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("published" + t.String())
			}
		}
	}()
	// Simple Synchronous Publisher
	time.Sleep(200 * time.Minute)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Duration(50) * time.Second)
	// Unsubscribe
	// Close connection
	sc.Close()
}
