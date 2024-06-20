package main

type AppConfig struct {
	Name              string
	Port              string
	HeartBeatInterval int
}

var Config = &AppConfig{
	Name:              "Load Balancer",
	Port:              ":8081",
	HeartBeatInterval: 1000,
}
