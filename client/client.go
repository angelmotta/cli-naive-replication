package client

import (
	"github.com/angelmotta/cli-naive-replication/internal/config"
	"github.com/angelmotta/cli-naive-replication/internal/exchangestore"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

type CmdLog struct {
	SendTime    time.Time     // the send time of this client command
	ReceiveTime time.Time     // the receive time of client command
	Duration    time.Duration // the calculated latency for this command (ReceiveTime - SendTime)
}

// ClientMetrics represents the performance metrics of a client
type ClientMetrics struct {
	MinLatency    float64
	MaxLatency    float64
	AvgLatency    float64
	P90Latency    float64
	P99Latency    float64
	Mid80Duration float64 // the middle 80% of the duration of the commands (seconds)
	Mid80Reqs     int     // Amount of requests in the middle 80% of the total requests
}

type Client struct {
	ClientId                   uint32
	servers                    []string
	exchangeStoreConn          []*exchangestore.ExchangeStore
	startWorkload, endWorkload time.Time
	RequestsExecuted           int
	CommandsLog                []CmdLog
	PerfMetrics                ClientMetrics
}

// New returns a new client
func New(idClient uint32, servers []string) *Client {
	n := len(servers)
	c := &Client{
		ClientId:          idClient,
		servers:           make([]string, n),
		exchangeStoreConn: make([]*exchangestore.ExchangeStore, n),
		CommandsLog:       make([]CmdLog, config.Global.NClientRequests),
	}
	copy(c.servers, servers)
	// Map each Redis servers to this client
	for i := 0; i < n; i++ {
		log.Printf("Mapping ClientId #%v with Redis server %v: %v", c.ClientId, i, c.servers[i])
		// Connect to ExchangeStore Replica 1
		eStore, err := exchangestore.New(c.servers[i])
		if err != nil {
			log.Panic("something happened New ExchangeStore: ", err)
		}
		c.exchangeStoreConn[i] = eStore
	}

	return c
}

func (c *Client) CloseLoopClient(wg *sync.WaitGroup, numSecs int) {
	defer wg.Done() // Decrement the counter when goroutine complete
	ClientTimeout := time.Duration(numSecs) * time.Second
	ticker := time.NewTicker(ClientTimeout) // channel to receive timeout
	log.Printf("ClientId #%v, started CloseLoop...", c.ClientId)
	c.startWorkload = time.Now()
MainLoopClient:
	for i := 0; i < config.Global.NClientRequests; i++ {
		select {
		case <-ticker.C:
			break MainLoopClient
		default:
			c.CommandsLog[i].SendTime = time.Now()
			c.sendOneRequest(i)
			c.CommandsLog[i].ReceiveTime = time.Now()
			c.CommandsLog[i].Duration = c.CommandsLog[i].ReceiveTime.Sub(c.CommandsLog[i].SendTime)
			c.RequestsExecuted += 1
		}
	}
	c.endWorkload = time.Now()
	log.Printf("ClientId #%v, finished CloseLoopClient!!", c.ClientId)
	elapsedSeconds := c.endWorkload.Sub(c.startWorkload).Seconds()
	log.Printf("ClientId #%v, ClooseLoop duration: %v seconds", c.ClientId, elapsedSeconds)
}

func (c *Client) getRandomCurrencyPrice() string {
	min := 3708
	max := 3910
	n := 8 // Length: 8 bytes
	// Generate Random
	valCurrency := float64(rand.Intn(max-min)+min) / 1000

	//valPrice := fmt.Sprintf("%f", valCurrency)
	valPrice := strconv.FormatFloat(valCurrency, 'f', 6, 64)
	if len(valPrice) != n {
		log.Panicf("got error creating random price value %v in ExchangeStore: this length is not %v bytes", valPrice, n)
	}
	return valPrice
}

func (c *Client) sendOneRequest(sn int) int {
	valPrice := c.getRandomCurrencyPrice()
	typeOp := rand.Intn(2) // typeOp: 0 is Set, 1 is Get
	if typeOp == 0 {       // typeOp is Set
		//log.Printf("ClientId #%v, OpNum #%v: set %v", c.ClientId, sn, valPrice)
		// Loop over list of replicas (redis servers) and write data to each one
		for _, exchangeStore := range c.exchangeStoreConn {
			err := exchangeStore.SetExchange("usd_pen_", valPrice)
			if err != nil {
				log.Panicf("got error Set value operation #%v in ExchangeStore: %v", sn, err)
			}
		}
	} else { // typeOp is Get
		//log.Printf("ClientId #%v, OpNum #%v, : get usd_pen_", c.ClientId, sn)
		// Loop over list of replicas (redis servers) and get data from each one
		for _, exchangeStore := range c.exchangeStoreConn {
			_, err := exchangeStore.GetExchange("usd_pen_")
			if err != nil {
				log.Panicf("got error Get value operation #%v in ExchangeStore: %v", sn, err)
			}
		}
	}
	return 0
}

func (c *Client) CalculateLatency() {
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
	cmdLogs := make([]CmdLog, lenCmdLogs)
	idx := 0
	for i := 0; i < len(c.CommandsLog); i++ {
		if i < int(float64(cmdsExecuted)*0.10) || i >= int(float64(cmdsExecuted)*0.90) {
			continue
		}
		if idx < lenCmdLogs {
			cmdLogs[idx] = c.CommandsLog[i]
			idx++
		} else {
			//log.Printf("cmdLogs array is ready with filtered data, idx: %v, lenCmdLogs: %v", idx, lenCmdLogs)
			break
		}
	}
	// Get duration time of middle 80% of commands log and amount of commands during this time
	mid80Start := cmdLogs[0].SendTime
	mid80End := cmdLogs[len(cmdLogs)-1].ReceiveTime
	mid80Duration := mid80End.Sub(mid80Start)
	c.PerfMetrics.Mid80Duration = mid80Duration.Seconds()
	c.PerfMetrics.Mid80Reqs = len(cmdLogs)

	// Sort commands log by latency (duration of request
	sort.Slice(cmdLogs, func(i, j int) bool {
		return cmdLogs[i].Duration < cmdLogs[j].Duration
	})

	var sumLatency float64
	for _, cmd := range cmdLogs {
		sumLatency += float64(cmd.Duration.Milliseconds())
	}

	c.PerfMetrics.MinLatency = float64(cmdLogs[0].Duration.Milliseconds())
	c.PerfMetrics.MaxLatency = float64(cmdLogs[len(cmdLogs)-1].Duration.Milliseconds())
	c.PerfMetrics.AvgLatency = sumLatency / float64(len(cmdLogs))
	c.PerfMetrics.P90Latency = float64(cmdLogs[int(float64(len(cmdLogs))*0.90)].Duration.Milliseconds())
	c.PerfMetrics.P99Latency = float64(cmdLogs[int(float64(len(cmdLogs))*0.99)].Duration.Milliseconds())

	log.Printf("ClientId #%v, Min Latency: %v milliseconds", c.ClientId, c.PerfMetrics.MinLatency)
	log.Printf("ClientId #%v, Max Latency: %v milliseconds", c.ClientId, c.PerfMetrics.MaxLatency)
	log.Printf("ClientId #%v, Avg Latency: %v milliseconds", c.ClientId, c.PerfMetrics.AvgLatency)
	log.Printf("ClientId #%v, P90 Latency: %v milliseconds", c.ClientId, c.PerfMetrics.P90Latency)
	log.Printf("ClientId #%v, P99 Latency: %v milliseconds", c.ClientId, c.PerfMetrics.P99Latency)
}
