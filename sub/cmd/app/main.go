package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	server "test/internal/http-server"
	model "test/internal/model"
	"time"

	"github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientID  = "order2"
	channel   = "order"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	s := server.NewServer()

	sc, err := stan.Connect(clusterID, clientID)

	if err != nil {
		log.Println("cant connect to nats")
	}

	sc.Subscribe(channel, func(msg *stan.Msg) {
		var order model.OrderMessage

		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println("error unmarshalling json")
		}

		err = order.Validate()
		if err != nil {
			fmt.Println(err)
		} else {
			s.Cache.Add(order)
		}

		err = s.DB.InsertMessage(order)
		if err != nil {
			fmt.Println(err)
		}
	})

	s.Start(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()
}
