package client

import (
	"context"
	"fmt"
	"github.com/angelmotta/cli-naive-replication/internal/exchangestore"
	"log"
	"math/rand"
	"sync"
)

type Client struct {
	Svr1Addr string
	Svr2Addr string
	replica1 *exchangestore.ExchangeStore
	replica2 *exchangestore.ExchangeStore
}

var ctx = context.Background()

// New returns a new client
func New(svr1Addr, svr2Addr string) *Client {
	c := &Client{
		Svr1Addr: svr1Addr,
		Svr2Addr: svr2Addr,
	}
	// Connect to ExchangeStore Replica 1
	r1, err := exchangestore.New("localhost:6380")
	if err != nil {
		log.Panic("something happened New ExchangeStore", err)
	}
	c.replica1 = r1
	// Connect to ExchangeStore Replica 2
	r2, err := exchangestore.New("localhost:6381")
	if err != nil {
		log.Panic("something happened New ExchangeStore", err)
	}
	c.replica2 = r2

	return c
}

func (c *Client) TestInsertions(priceExchange float64, wg *sync.WaitGroup, n int) {
	defer wg.Done() // Decrement the counter when goroutine complete

	log.Println("TestInsertions execution started...")
	min := 37100
	max := 39100
	for i := 0; i < n; i++ {
		// Generate Random
		valCurrency := float64(rand.Intn(max-min)+min) / 10000
		valPrice := fmt.Sprintf("%f", valCurrency)
		// Writes to Replica1
		err := c.replica1.SetExchange("sol-dollar", valPrice)
		if err != nil {
			log.Panicf("got error Set value %v in ExchangeStore: %v", priceExchange, err)
		}
		// Writes to replica2
		err = c.replica2.SetExchange("sol-dollar", valPrice)
		if err != nil {
			log.Panicf("got error Set value %v in ExchangeStore: %v", priceExchange, err)
		}
		priceExchange = priceExchange + 0.002
	}
	log.Println("TestInsertions execution done")
}
