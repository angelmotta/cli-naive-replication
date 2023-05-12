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
	ClientId uint32
	Svr1Addr string
	Svr2Addr string
	Svr3Addr string
	replica1 *exchangestore.ExchangeStore
	replica2 *exchangestore.ExchangeStore
	replica3 *exchangestore.ExchangeStore
}

var ctx = context.Background()

// New returns a new client
func New(idClient uint32, svr1Addr, svr2Addr, svr3Addr string) *Client {
	c := &Client{
		ClientId: idClient,
		Svr1Addr: svr1Addr,
		Svr2Addr: svr2Addr,
		Svr3Addr: svr3Addr,
	}
	// Connect to ExchangeStore Replica 1
	r1, err := exchangestore.New(svr1Addr)
	if err != nil {
		log.Panic("something happened New ExchangeStore", err)
	}
	c.replica1 = r1
	// Connect to ExchangeStore Replica 2
	r2, err := exchangestore.New(svr2Addr)
	if err != nil {
		log.Panic("something happened New ExchangeStore", err)
	}
	c.replica2 = r2

	// Connect to ExchangeStore Replica 3
	r3, err := exchangestore.New(svr3Addr)
	if err != nil {
		log.Panic("something happened New ExchangeStore", err)
	}
	c.replica3 = r3

	return c
}

func (c *Client) CloseLoopClient(wg *sync.WaitGroup, n int) {
	defer wg.Done() // Decrement the counter when goroutine complete

	log.Println("CloseLoopClient execution started...")
	for i := 0; i < n; i++ {
		//log.Printf("iteration #%v from clientId #%v", i, c.ClientId)
		valPrice := c.GetRandomCurrencyPrice()
		typeOp := rand.Intn(2) // typeOp: 0 is Set, 1 is Get
		if typeOp == 0 {
			log.Printf("ClientId #%v, OpNum #%v: set %v", c.ClientId, i, valPrice)
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

			// Writes to replica3
			err = c.replica3.SetExchange("usd_pen_", valPrice)
			if err != nil {
				log.Panicf("got error Set value %v in ExchangeStore: %v", i, err)
			}
		} else { // typeOp is Get
			log.Printf("ClientId #%v, OpNum #%v, : get usd_pen_", c.ClientId, i)
			// Reads from Replica1
			_, err := c.replica1.GetExchange("usd_pen_")
			if err != nil {
				log.Panicf("got error Get value operation #%v in ExchangeStore: %v", i, err)
			}

			// Reads from Replica2
			_, err = c.replica2.GetExchange("usd_pen_")
			if err != nil {
				log.Panicf("got error Get value operation #%v in ExchangeStore: %v", i, err)
			}

			// Reads from Replica3
			_, err = c.replica3.GetExchange("usd_pen_")
			if err != nil {
				log.Panicf("got error Get value operation #%v in ExchangeStore: %v", i, err)
			}
		}
	}
	log.Println("CloseLoopClient execution done")
}

func (c *Client) GetRandomCurrencyPrice() string {
	min := 3708
	max := 3910
	n := 8 // Length: 8 bytes
	// Generate Random
	valCurrency := float64(rand.Intn(max-min)+min) / 1000
	//log.Printf("float currency: %v", valCurrency)

	//valPrice := fmt.Sprintf("%f", valCurrency)
	valPrice := strconv.FormatFloat(valCurrency, 'f', 6, 64)
	//log.Printf("my string currency: %v", valPrice)
	if len(valPrice) != n {
		log.Panicf("got error creating random price value %v in ExchangeStore: this length is not %v bytes", valPrice, n)
	}
	return valPrice
}
