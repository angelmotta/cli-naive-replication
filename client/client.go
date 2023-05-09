package client

import (
	"context"
	"github.com/angelmotta/cli-naive-replication/internal/exchangestore"
	"log"
	"math/rand"
	"strconv"
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

func (c *Client) TestInsertions(wg *sync.WaitGroup, n int) {
	defer wg.Done() // Decrement the counter when goroutine complete

	log.Println("TestInsertions execution started...")
	for i := 0; i < n; i++ {
		valPrice := c.GetRandomCurrencyPrice()
		// Writes to Replica1
		err := c.replica1.SetExchange("usd_pen_", valPrice)
		if err != nil {
			log.Panicf("got error Set value opeartion #%v in ExchangeStore: %v", i, err)
		}
		// Writes to replica2
		err = c.replica2.SetExchange("usd_pen_", valPrice)
		if err != nil {
			log.Panicf("got error Set value %v in ExchangeStore: %v", i, err)
		}
	}
	log.Println("TestInsertions execution done")
}

func (c *Client) GetRandomCurrencyPrice() string {
	min := 3708
	max := 3910
	n := 8 // Length: 8 bytes
	// Generate Random
	valCurrency := float64(rand.Intn(max-min)+min) / 1000
	log.Printf("float currency: %v", valCurrency)
	//valPrice := fmt.Sprintf("%f", valCurrency)
	valPrice := strconv.FormatFloat(valCurrency, 'f', 6, 64)
	log.Printf("my string currency: %v", valPrice)
	if len(valPrice) != n {
		log.Panicf("got error creating random price value %v in ExchangeStore: this length is not %v bytes", valPrice, n)
	}
	return valPrice
}
