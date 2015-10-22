package main

import (
	"fmt"
	"gopkg.in/redis.v3"

	"bufio"
	"flag"
	"os"
	"strings"
)

var msgchan chan *redis.Message

var (
	addr     string // redis host addr
	password string // redis password
	db       int64  // redis db
)

func main() {
	// flags
	flag.StringVar(&addr, "addr", "0.0.0.0:6379", "-addr=0.0.0.0:6379")
	flag.StringVar(&password, "password", "", "-password=xxx")
	flag.Int64Var(&db, "db", 0, "-db=n")
	flag.Parse()

	inputchan := make(chan string, 1)
	msgchan = make(chan *redis.Message, 10)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	fmt.Printf("Ping Type : %T \n", client.Ping)

	// clientErrHelper(client.Ping)

	go func(msgchan chan *redis.Message) {
		foosub, _ := client.Subscribe("foo")
		fmt.Println("foosub : ", foosub.String())
		for {
			msg, _ := foosub.ReceiveMessage() //
			msgchan <- msg
		}

	}(msgchan)

	go func(msgchan chan *redis.Message) {
		booksub, _ := client.Subscribe("book")
		for {
			msg, _ := booksub.ReceiveMessage() //
			msgchan <- msg
		}

	}(msgchan)

	// publish
	go func(inputchan chan string) {
		r := bufio.NewReader(os.Stdin)
		for {
			msg, _ := r.ReadSlice('\n')
			if string(msg[:len("PUBLISH")]) == "PUBLISH" {
				inputchan <- string(msg[len("PUBLISH")+1:])
			}

			// if string(msg[:len("UNSUBSCRIBE")]) == "UNSUBSCRIBE" {
			// 	client.Unsub string(msg[len("UNSUBSCRIBE")+1:])
			// }
		}
	}(inputchan)
	channels := client.PubSubChannels("*")
	fmt.Println("channels : ", channels.Val())

	for {

		select {
		case recMsg := <-msgchan:
			fmt.Println("interested : ")
			fmt.Println(recMsg)
		case inputMsg := <-inputchan:
			pubMsg := strings.Fields(inputMsg)
			r := client.Publish(pubMsg[0], pubMsg[1])
			fmt.Println("publish r: ", r)

		}
	}

}

// func clientErrHelper(f interface{}) error {
// 	switch function := f.(type) {
// 	case func() *redis.StatusCmd:
// 		str, err := function().Result()
// 		if err != nil {
// 			fmt.Println("result error : ", err.Error())
// 			return err
// 		}
// 		fmt.Println("result : ", str)
// 		return nil
// 	}
// 	return nil
// }
