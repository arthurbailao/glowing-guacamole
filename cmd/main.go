package main

import (
	"fmt"
	"os"

	"github.com/arthurbailao/cint-ad/cmd/consumer"
	"github.com/arthurbailao/cint-ad/cmd/producer"
)

func main() {
	if len(os.Args) != 2 {
		fail()
	}

	switch os.Args[1] {
	case "producer":
		producer.Run()

	case "consumer":
		consumer.Run()

	default:
		fail()
	}
}

func fail() {
	fmt.Println("Usage: ./cmd <command>")
	fmt.Println("The comands are: producer, consumer")
	os.Exit(2)
}
