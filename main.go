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
	log.Println("*** Client test replication finished ***\n")
	log.Printf("Total Elapsed time Test Replication: %v seconds\n", elapsedSeconds)

	// Generate Performance metrics by client
	generateLatencyMetrics(clients)
	// Print summary results per client
	printSummaryResults(clients)
}

func generateLatencyMetrics(clients []*client.Client) {
	log.Println("--- Performance metrics by client ---")
	// Calculate the average latency per client
	for i := 0; i < config.Global.NClients; i++ {
		log.Printf("** Performance metrics ClientId: %v **", i)
		clients[i].CalculateLatency()
	}
}

func printSummaryResults(clients []*client.Client) {
	generalThroughtput := 0.0
	mid80Throughput := 0.0
	mid80MinLat := clients[0].PerfMetrics.MinLatency
	mid80MaxLat := clients[0].PerfMetrics.MaxLatency
	mid80AvgLat := 0.0
	mid80P90Lat := 0.0
	mid80P99Lat := 0.0
	log.Println("--- Throughput Results by Client ---")
	for i := 0; i < len(clients); i++ {
		// Throughput per client
		log.Printf("--- Results ClientId #%v ---", clients[i].ClientId)
		log.Printf("ClientId #%v, Total requests executed: %v", clients[i].ClientId, clients[i].RequestsExecuted)
		throughput := math.Round(float64(clients[i].RequestsExecuted) / float64(config.Global.DurationTest))
		generalThroughtput += throughput
		// Mid 80% metrics per client
		clientMid80thr := math.Round(float64(clients[i].PerfMetrics.Mid80Reqs) / clients[i].PerfMetrics.Mid80Duration)
		mid80Throughput += clientMid80thr
		// Min Latency of system
		if clients[i].PerfMetrics.MinLatency < mid80MinLat {
			mid80MinLat = clients[i].PerfMetrics.MinLatency
		}
		// Max Latency of system
		if clients[i].PerfMetrics.MaxLatency > mid80MaxLat {
			mid80MaxLat = clients[i].PerfMetrics.MaxLatency
		}
		// Latencies of system
		mid80AvgLat += clients[i].PerfMetrics.AvgLatency
		mid80P90Lat += clients[i].PerfMetrics.P90Latency
		mid80P99Lat += clients[i].PerfMetrics.P99Latency
		log.Printf("ClientId #%v, Throughput: %v req/sec", clients[i].ClientId, throughput)
		log.Printf("ClientId #%v, Mid80Throughput: %v req/sec", clients[i].ClientId, clientMid80thr)
	}
	log.Println("--- General Performance Server Results ---")
	log.Printf("General Throughput: %v req/sec", generalThroughtput)
	log.Printf("Throughput (Mid80): %v req/sec", mid80Throughput)
	log.Printf("Min Latency (Mid80): %v ms", mid80MinLat)
	log.Printf("Max Latency (Mid80): %v ms", mid80MaxLat)
	log.Printf("Avg Latency (Mid80): %v ms", mid80AvgLat/float64(len(clients)))
	log.Printf("P90 Latency (Mid80): %v ms", mid80P90Lat/float64(len(clients)))
	log.Printf("P99 Latency (Mid80): %v ms", mid80P99Lat/float64(len(clients)))
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
