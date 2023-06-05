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
	c.NClients = 15              // number of concurrent clients
	c.DurationTest = 20          // seconds
	c.RedisList = []string{"172.31.29.50:6379", "172.31.22.226:6380", "172.31.22.56:6381"}
}
