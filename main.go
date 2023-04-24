package main

import (
	"fmt"
	"github.com/angelmotta/cli-naive-replication/internal/exchangestore"
	"log"
)

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

func main() {
	log.Println("*** client naive replication started... ***")
	r1, err := exchangestore.New("localhost:6381")
	if err != nil {
		log.Panic("something happened New ExchangeStore", err)
	}

	testClient(r1, 3.810)
	log.Println(r1.GetExchange("sol-dollar"))
	log.Println("*** client naive replication started... ***")
}
