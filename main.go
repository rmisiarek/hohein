package main

import "log"

func main() {
	// client := getClient(15)

	// p := Payload{
	// 	url:    "https://www.golang.org",
	// 	method: "GET",
	// 	key:    "x-forwarded-for",
	// 	value:  "0177.0.0.01",
	// }

	opts := &HohinOptions{
		pathHosts:   "",
		pathHeaders: "headers.txt",
		pathValues:  "values.txt",
		output:      "",
		timeout:     15,
	}

	_, err := NewHohin(opts)
	if err != nil {
		log.Fatalln(err)
	}

	// resp, err := client.request(p, true)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Printf("%#v\n", resp)
}
