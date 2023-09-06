package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	stan "github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientID  = "order"
	channel   = "order"
)

func Init() stan.Conn {
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Println("cant connect to stan")
	}
	return sc
}

func ReadFilePublish(sc stan.Conn) {
	filepth, err := os.Executable()
	reader := bufio.NewReader(os.Stdin)

	if err != nil {
		return
	}
	for {
		fmt.Println("Enter filename in ./jsons folder")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\r\n\"")
		if text == "" {
			fmt.Println("Goodbye")
			return
		}

		filepath := path.Join(path.Dir(filepth), "jsons", text)
		fmt.Println(filepath)

		file, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Println("Cannot read a file, try another")
			continue
		}

		err = sc.Publish(channel, file)
		if err != nil {
			fmt.Println("Cannot publish message, please retry")
			continue
		}
	}
}

func main() {
	sc := Init()
	ReadFilePublish(sc)
}
