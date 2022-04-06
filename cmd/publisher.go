package main

import (
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
)

func main() {
	sc, err := stan.Connect("test-cluster", "pub")
	if err != nil {
		log.Fatal(err)
	}

	jsonFile, err := ioutil.ReadFile("../model.json")
	if err != nil {
		log.Fatal(err)
	}

	sc.Publish("foo", jsonFile)
	sc.Close()
}
