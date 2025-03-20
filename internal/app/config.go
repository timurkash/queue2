package app

import (
	"flag"
	"time"
)

type Config struct {
	Port           int
	DefaultTimeout time.Duration
	MaxQueues      int
	QueueCapacity  int
}

func ParseFlags() Config {
	var cfg Config
	flag.IntVar(&cfg.Port, "port", 8080, "queue2 port")
	flag.DurationVar(&cfg.DefaultTimeout, "timeout", 30*time.Second, "default timeout")
	flag.IntVar(&cfg.MaxQueues, "max-queues", 100, "max queues count")
	flag.IntVar(&cfg.QueueCapacity, "queue-capacity", 1000, "max messages per queue")
	flag.Parse()
	return cfg
}
