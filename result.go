package main

type Result struct {
	payloads              Payload
	statusCode            int
	headerValue           string
	host                  string
	location              string
	reflectedKeys         []string
	reflectedValues       []string
	reflectedValuesInBody []string
	confirmed             bool
}

type ResultWrapper struct {
	result Result
	err    error
}

func handleResult(r ResultWrapper) {
	if r.err != nil {
		red.Printf("\t==> %s\n", r.err.Error())
		return
	}

	if r.result.confirmed {
		greenBold.Printf("\t==> status code: %d | payload: %s | CONFIRMED!\n", r.result.statusCode, r.result.headerValue)
	} else {
		blue.Printf("\t==> status code: %d | payload: %s\n", r.result.statusCode, r.result.headerValue)
	}

	if len(r.result.reflectedKeys) != 0 {
		greenBold.Printf("\t\t>> found reflected header key:\n")
		for _, v := range r.result.reflectedKeys {
			greenBold.Printf("\t\t\t>> %s\n", v)
		}
	}

	if len(r.result.reflectedValues) != 0 {
		greenBold.Printf("\t\t>> found reflected header value:\n")
		for _, v := range r.result.reflectedValues {
			greenBold.Printf("\t\t\t>> %s\n", v)
		}
	}

	if len(r.result.reflectedValuesInBody) != 0 {
		greenBold.Printf("\t\t>> found reflected header value in body:\n")
		for _, v := range r.result.reflectedValuesInBody {
			greenBold.Printf("\t\t\t>> %s\n", v)
		}
	}
}
