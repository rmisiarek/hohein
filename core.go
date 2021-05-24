package main

type HohinResult struct {
	payload            Payload
	responseStatusCode int
	responseURL        string
	reflectedKey       string
	reflectedValue     string
	reflectedValueBody string
	location           string
	confirmed          bool
}

type Payload struct {
	url    string
	method string
	key    string
	value  string
}
