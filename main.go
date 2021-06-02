package main

import (
	"log"
)

func main() {
	opts := &HohinOptions{
		pathHosts:   "hosts.txt",
		pathHeaders: "headers.txt",
		pathValues:  "values.txt",
		output:      "",
		timeout:     15,
	}

	h, err := NewHohin(opts)
	if err != nil {
		log.Fatalln(err)
	}

	h.Start()
}
