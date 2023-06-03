package main

import (
	"fmt"
	"github.com/angelmotta/cli-naive-replication/client"
	"github.com/angelmotta/cli-naive-replication/internal/config"
	"github.com/angelmotta/cli-naive-replication/internal/exchangestore"
	"log"
	"math"
	"sync"
	"time"
)

func main() {
	log.Println("*** Client test replication started ***")
	// General configuration
	config.Global.Init()

	// Create clients
	clients := make([]*client.Client, config.Global.NClients)
	for i := 0; i < config.Global.NClients; i++ {
		clients[i] = client.New(uint32(i), config.Global.RedisList)
	}

	// Run concurrently clients
	startTime := time.Now()
	wg := new(sync.WaitGroup)
	for i := 0; i < config.Global.NClients; i++ {
		wg.Add(1)
		go clients[i].CloseLoopClient(wg, config.Global.DurationTest)
	}

	// Wait until clients finish their workloads
	log.Println("waiting until clients finish their workloads...")
	wg.Wait()
	endTime := time.Now()

	elapsedSeconds := endTime.Sub(startTime).Seconds()
	log.Println("*** Client test replication finished ***")
	log.Printf("Elapsed time in Test Replication: %v seconds\n", elapsedSeconds)

	// Print summary results per client
	printSummaryResults(clients)
}

func printSummaryResults(clients []*client.Client) {
	generalThroughtput := 0.0
	log.Println("--- Summary results per client ---")
	for i := 0; i < config.Global.NClients; i++ {
		log.Printf("ClientId #%v, requests executed: %v", clients[i].ClientId, clients[i].RequestsExecuted)
		// print throughput per client
		throughput := math.Round(float64(clients[i].RequestsExecuted) / float64(config.Global.DurationTest))
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
