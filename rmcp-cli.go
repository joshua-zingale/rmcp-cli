package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/joshua-zingale/rmcp-cl/internal"
)

func main() {

	host := flag.String("host", "", "The host address (e.g. '127.0.0.1') for the remote mcp host")
	port := flag.String("port", "", "The port (e.g. '80' or '5000') of the remote mcp host")

	flag.Parse()

	if *host == "" {
		*host = "http://127.0.0.1"
	}
	if *port == "" {
		*port = "80"
	}

	client := internal.NewClient(*host, *port)

	r, err := client.ListServers()
	if err != nil {
		panic(err)
	}

	fmt.Println(internal.FormatServerList(r))

	mux := internal.NewCommandMux(os.Stdout, &client)

	mux.AddHandlerFunc(`^/.+`, internal.HandleCommand)
	mux.AddHandlerFunc("(.+)", internal.HandleUserMessage)

	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		err := mux.Handle(ctx, line)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}
