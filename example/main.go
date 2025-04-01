package main

import (
	"fmt"
	"log"

	picoclient "github.com/stipochka/pico_go/client"
)

func main() {
	var filename string
	fmt.Println("enter slave path")
	fmt.Scan(&filename)
	client, err := picoclient.NewClient(filename)
	defer client.Close()
	if err != nil {
		log.Fatalf("failed to create client")
	}

	resp, err := client.SetReadingPeriodRequest("temp0", 10)
	if err != nil {
		log.Fatal("error with heartbit_request", err)
	}
	fmt.Println(resp)
	fmt.Println(string(resp.Buffer))
}
