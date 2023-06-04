package main

import (
	"fmt"
	"github.com/angelmotta/cli-naive-replication/client"
	"github.com/angelmotta/cli-naive-replication/internal/config"
	"github.com/angelmotta/cli-naive-replication/internal/exchangestore"
	"log"
	"math"
	"sort"
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
	// Generate Performance metrics
	generateLatencyMetrics(clients)
}

func generateLatencyMetrics(clients []*client.Client) {
	log.Println("--- Performance metrics ---")
	// Calculate the average latency per client
	for i := 0; i < config.Global.NClients; i++ {
		calculateLatency(clients[i])
	}
}

func calculateLatency(c *client.Client) {
	// Validate number of commands executed
	var cmdsExecuted int
	for i := 0; i < len(c.CommandsLog); i++ {
		if c.CommandsLog[i].Duration == time.Duration(0) {
			break
		}
		cmdsExecuted++
	}

	if cmdsExecuted != c.RequestsExecuted {
		log.Panicf("something happened, cmdsExecuted: %v, RequestsExecuted: %v", cmdsExecuted, c.RequestsExecuted)
	}

	// Exclude head and tails requests (20%) from the commands log
	lenCmdLogs := int(float64(cmdsExecuted) * 0.80)
	cmdLogs := make([]client.CmdLog, lenCmdLogs)
	idx := 0
	for i := 0; i < len(c.CommandsLog); i++ {
		if i < int(float64(cmdsExecuted)*0.10) || i >= int(float64(cmdsExecuted)*0.90) {
			continue
		}
		if idx < lenCmdLogs {
			cmdLogs[idx] = c.CommandsLog[i]
			idx++
		} else {
			log.Printf("cmdLogs array is ready with filtered data, idx: %v, lenCmdLogs: %v", idx, lenCmdLogs)
			break
		}
	}

	// Calculate latencies
	sort.Slice(cmdLogs, func(i, j int) bool {
		return cmdLogs[i].Duration < cmdLogs[j].Duration
	})

	var sumLatency float64
	for _, cmd := range cmdLogs {
		sumLatency += float64(cmd.Duration.Milliseconds())
	}

	minLat := float64(cmdLogs[0].Duration.Milliseconds())
	maxLat := float64(cmdLogs[len(cmdLogs)-1].Duration.Milliseconds())
	avgLat := sumLatency / float64(len(cmdLogs))
	p90Lat := float64(cmdLogs[int(float64(len(cmdLogs))*0.90)].Duration.Milliseconds())
	p99Lat := float64(cmdLogs[int(float64(len(cmdLogs))*0.99)].Duration.Milliseconds())
	log.Printf("ClientId #%v, Min Latency: %v milliseconds", c.ClientId, minLat)
	log.Printf("ClientId #%v, Max Latency: %v milliseconds", c.ClientId, maxLat)
	log.Printf("ClientId #%v, Avg Latency: %v milliseconds", c.ClientId, avgLat)
	log.Printf("ClientId #%v, P90 Latency: %v milliseconds", c.ClientId, p90Lat)
	log.Printf("ClientId #%v, P99 Latency: %v milliseconds", c.ClientId, p99Lat)
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
