package config

var Global Config

type Config struct {
	NClientRequests int
	NClients        int
	DurationTest    int
	RedisList       []string
}

func (c *Config) Init() {
	c.NClientRequests = 10000000 // the default value
	c.NClients = 1               // number of concurrent clients
	c.DurationTest = 5           // seconds
	c.RedisList = []string{"localhost:6380", "localhost:6381", "localhost:6382"}
}
