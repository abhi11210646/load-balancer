package main

type AppConfig struct {
	Name              string
	Port              string
	Servers           []string
	HeartBeatInterval int
}

var Config = AppConfig{
	Name:              "Load Balancer",
	Port:              ":3000",
	Servers:           []string{},
	HeartBeatInterval: 1000,
}
