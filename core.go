package main

type HohinResult struct {
	statusCode     int
	url            string
	headerKey      string
	headerValue    string
	keyReflected   []string
	valueReflected []string
	location       string
	redirect       bool
	confirmed      bool
}

type HohinPayload struct {
	url         string
	method      string
	headerKey   string
	headerValue string
}
