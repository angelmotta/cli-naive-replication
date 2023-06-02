package main

import (
	"fmt"
	"github.com/angelmotta/cli-naive-replication/client"
	"github.com/angelmotta/cli-naive-replication/internal/exchangestore"
	"log"
	"math"
	"sync"
	"time"
)

func main() {
	log.Println("*** Client test replication started ***")
	// General configuration
	RedisList := []string{"localhost:6380", "localhost:6381", "localhost:6382"}
	NClients := 6
	DurationTest := 10 // seconds

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
		go clients[i].CloseLoopClient(wg, DurationTest)
	}

	// Wait until both clients finish their workloads
	log.Println("waiting to finish both clients")
	wg.Wait()
	endTime := time.Now()

	elapsedSeconds := endTime.Sub(startTime).Seconds()
	log.Println("*** Client test replication finished ***")
	log.Printf("Elapsed time in Test Replication: %v seconds\n", elapsedSeconds)

	// Print summary results per client
	generalThroughtput := 0.0
	log.Println("--- Summary results per client ---")
	for i := 0; i < NClients; i++ {
		log.Printf("ClientId #%v, requests executed: %v", clients[i].ClientId, clients[i].RequestsExecuted)
		// print throughput per client
		throughput := math.Round(float64(clients[i].RequestsExecuted) / elapsedSeconds)
		generalThroughtput += throughput
		log.Printf("ClientId #%v, Throughput: %v req/sec", clients[i].ClientId, throughput)
	}
	log.Println("--- Summary General Results ---")
	log.Printf("General Throughput: %v req/sec", generalThroughtput)
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
