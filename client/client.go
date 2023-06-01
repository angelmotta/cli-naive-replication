package client

import (
	"github.com/angelmotta/cli-naive-replication/internal/exchangestore"
	"log"
	"math/rand"
	"strconv"
	"sync"
)

type Client struct {
	ClientId          uint32
	servers           []string
	exchangeStoreConn []*exchangestore.ExchangeStore
}

// New returns a new client
func New(idClient uint32, servers []string) *Client {
	n := len(servers)
	c := &Client{
		ClientId:          idClient,
		servers:           make([]string, n),
		exchangeStoreConn: make([]*exchangestore.ExchangeStore, n),
	}
	copy(c.servers, servers)
	// Map each Redis servers to this client
	for i := 0; i < n; i++ {
		log.Printf("Mapping ClientId #%v with Redis server %v: %v", c.ClientId, i, c.servers[i])
		// Connect to ExchangeStore Replica 1
		eStore, err := exchangestore.New(c.servers[i])
		if err != nil {
			log.Panic("something happened New ExchangeStore", err)
		}
		c.exchangeStoreConn[i] = eStore
	}

	return c
}

func (c *Client) CloseLoopClient(wg *sync.WaitGroup, n int) {
	defer wg.Done() // Decrement the counter when goroutine complete

	log.Printf("ClientId #%v, started CloseLoop...", c.ClientId)
	for i := 0; i < n; i++ {
		c.sendOneRequest(i)
	}
	log.Printf("ClientId #%v, finished CloseLoopClient!!", c.ClientId)
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

func (c *Client) sendOneRequest(sn int) {
	valPrice := c.getRandomCurrencyPrice()
	typeOp := rand.Intn(2) // typeOp: 0 is Set, 1 is Get
	if typeOp == 0 {       // typeOp is Set
		log.Printf("ClientId #%v, OpNum #%v: set %v", c.ClientId, sn, valPrice)
		// Loop redis servers (replicas) and write to each one
		for _, exchangeStore := range c.exchangeStoreConn {
			err := exchangeStore.SetExchange("usd_pen_", valPrice)
			if err != nil {
				log.Panicf("got error Set value operation #%v in ExchangeStore: %v", sn, err)
			}
		}
	} else { // typeOp is Get
		log.Printf("ClientId #%v, OpNum #%v, : get usd_pen_", c.ClientId, sn)
		// Loop redis servers (replicas) and write to each one
		for _, exchangeStore := range c.exchangeStoreConn {
			_, err := exchangeStore.GetExchange("usd_pen_")
			if err != nil {
				log.Panicf("got error Get value operation #%v in ExchangeStore: %v", sn, err)
			}
		}
	}
}
