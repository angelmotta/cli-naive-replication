package main

import (
	"fmt"
	"github.com/angelmotta/cli-naive-replication/client"
	"github.com/angelmotta/cli-naive-replication/internal/exchangestore"
	"log"
	"sync"
	"time"
)

func main() {
	//initialTestApproach();
	log.Println("*** Client test replication started ***")
	// Parameters configuration
	RedisList := []string{"localhost:6380", "localhost:6381", "localhost:6382"}
	NClients := 2
	NReqs := 50

	// Create clients
	clients := make([]*client.Client, NClients)
	for i := 0; i < NClients; i++ {
		clients[i] = client.New(uint32(i), RedisList)
	}

	// Run concurrently clients
	startTime := time.Now()
	wg := new(sync.WaitGroup)
	for i := 0; i < NClients; i++ {
		wg.Add(1)
		go clients[i].CloseLoopClient(wg, NReqs)
	}

	// Wait until both clients finish their workloads
	log.Println("waiting to finish both clients")
	wg.Wait()
	endTime := time.Now()

	// Print results
	elapsedSeconds := endTime.Sub(startTime).Seconds()
	log.Println("*** Client test replication finished ***")
	log.Printf("Elapsed time in Test Replication: %v seconds\n", elapsedSeconds)
}

// initialTestApproach was the initial old approach
func initialTestApproach() {
	log.Println("*** client naive replication started... ***")
	r1, err := exchangestore.New("localhost:6381")
	if err != nil {
		log.Panic("something happened New ExchangeStore", err)
	}
	testClient(r1, 3.810)
	log.Println(r1.GetExchange("sol-dollar"))
	log.Println("*** client naive replication started... ***")
}

func testClient(conn *exchangestore.ExchangeStore, priceExchange float64) {
	log.Println("start testClient execution")
	for i := 0; i < 11; i++ {
		valPrice := fmt.Sprintf("%f", priceExchange)
		err := conn.SetExchange("sol-dollar", valPrice)
		if err != nil {
			log.Panicf("got error Set value %v in ExchangeStore: %v", priceExchange, err)
		}
		priceExchange = priceExchange + 0.002
	}
}
