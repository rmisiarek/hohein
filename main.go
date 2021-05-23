package main

import (
	"fmt"
	"log"
)

func main() {
	client := getClient(15)

	// client.referenceRequest("https://www.hackerone.com", "GET")

	resp, err := client.request("https://www.golang.org", "GET", "x-forwarded-for", "0177.0.0.01")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v\n", resp)
}
