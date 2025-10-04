package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

func main() {

	host := flag.String("host", "", "The host address (e.g. '127.0.0.1') for the remote mcp host")
	port := flag.String("port", "", "The port of the remote mcp host")

	flag.Parse()

	if *host == "" {
		*host = "http://127.0.0.1"
	}
	if *port == "" {
		*port = "80"
	}

	address := *host + ":" + *port

	req, err := http.NewRequest("GET", address+"/servers", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

}
