package exchangestore

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
)

type ExchangeStore struct {
	redis *redis.Client
}

var ctx = context.Background()

func New(addr string) (*ExchangeStore, error) {
	if addr == "" {
		log.Println("addr is blank")
		return nil, errors.New("DB Addr can not be blank")
	}
	rdbClient := redis.NewClient(&redis.Options{
		Addr:     addr, // Addr: "localhost:6379"
		Password: "",
		DB:       0,
	})

	// Test connection
	err := rdbClient.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &ExchangeStore{
		redis: rdbClient,
	}, nil
}

// GetExchange retrieves a Currency Exchange value from the Store layer
func (db *ExchangeStore) GetExchange(key string) (string, error) {
	val, err := db.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("key does not exist")
		return "", nil
	} else if err != nil {
		log.Println("got error Get value from Redis", err)
		return "", err
	}
	return val, err
}

// SetExchange set a Currency Exchange value to the Store layer
func (db *ExchangeStore) SetExchange(key string, val string) error {
	err := db.redis.Set(ctx, key, val, 0).Err()
	return err
}
